---
title: Java and Native Code in OpenJDK 13
layout: post
date: 2019-06-11T12:00:00.00Z
publish: true
---
# Contents

* [Intro](#intro)
* [Hands-on](#hands-on)
  * [The libCurl App](#libcurl)
  * [The libSSH2 App](#libssh)
* [Final Thoughts](#outro)
* [References](#refs)

# <a name="intro"></a>Intro

Last week I watched a pretty cool GOTO Conference talk [[1]](#refs) with Mikael Vidstedt, director of Software Engineering at Oracle, where he presents many of the upcoming and currently being worked on Java features. One of those, **Project Panama**, stood out as particularly interesting to me.

<blockquote class="twitter-tweet"><p lang="en" dir="ltr"><a href="https://twitter.com/hashtag/Java?src=hash&amp;ref_src=twsrc%5Etfw">#Java</a> has changed more in the last year than in the previous five years.<br><br>Learn about these changes and what Java&#39;s future holds in this <a href="https://twitter.com/hashtag/GOTOchgo?src=hash&amp;ref_src=twsrc%5Etfw">#GOTOchgo</a> talk with <a href="https://twitter.com/gsaab?ref_src=twsrc%5Etfw">@gsaab</a> and <a href="https://twitter.com/MikaelVidstedt?ref_src=twsrc%5Etfw">@MikaelVidstedt</a><a href="https://t.co/QcX8MnEQZL">https://t.co/QcX8MnEQZL</a></p>&mdash; GOTO (@GOTOcon) <a href="https://twitter.com/GOTOcon/status/1136248991254028289?ref_src=twsrc%5Etfw">June 5, 2019</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>

*Disclaimer: Project Panama is still in its early stages of development, so the contents of this article may not reflect its current state.*

Project Panama, available in early-access JDK 13 builds, is meant as a bridge between Java and native code. How is this useful? There are certain use cases where a piece of functionality is available as a library written in C or C++ and the requirement is to integrate it into a Java application. The current solution is to use the Java Native Interface (JNI) [[2]](#refs) which could require a vast amount of work and is also quite error-prone. You need to have a good grasp of the native library's insides and then write all the required implementations to bridge the Java interface with the native code. From my experience, calling native library functions that may change at some later point in time and also managing heap memory allocations could be a real challenge with JNI. To quote Mikael as he says it on the video:

> How many people here have written JNI or, you know, done JNI? We're so sorry for you!

This is where Panama comes into play. It provides new tools and API to simplify the bridging between Java and native code. It basically boils down to the following steps:

1. Generate Java bindings (interfaces) from existing C header file(s) using the **jextract** tool.
2. Invoke C functions via the extracted Java interface using the **java.foreign** API.

This allows you to concentrate on writing the core logic of your application, rather than fiddling with glue code and integration details.

# <a name="hands-on"></a>Hands-on

Project Panama's documentation pages [[3]](#refs) already provide a solid number of examples to start with. I'll just take a quick peek at how to bridge and run a *libCurl* Java app and then I'd like to present a more detailed example - a simple SSH client that I wrote based on *libSSH2*. 

I'm running these examples on a macOS, but with a few tweaks you should also be able to run them on a Linux installation.

## <a name="libcurl"></a>The libCurl App

How to download a web page using the native Curl library's implementation? Well, first we need to get and extract a Panama [OpenJDK build](http://jdk.java.net/panama/) archive.

Let's open up a shell and set the `JAVA_HOME` environment variable to where the OpenJDK build archive is extracted.

```bash
export JAVA_HOME=/opt/jdk-13.jdk
```

Now we need to generate the Java interfaces, the glue code that will bind Java code to the native library. This will produce a `curl.jar` file:

```bash
$JAVA_HOME/bin/jextract -t org.unix -L /usr/lib -lcurl \
  --record-library-path /usr/include/curl/curl.h \
  -o curl.jar
```

When we inspect the JAR file, we can see all the Curl API calls, as well as dependency bindings. The Curl API is available through the new **java.foreign** Java API.

![curl.jar](https://thepracticaldev.s3.amazonaws.com/i/5j56upyw7yz643tiwefx.png)

Now for a quick example. Here's a Java piece of code that fetches a web page and displays its contents on screen.

```java
import java.lang.invoke.*;
import java.foreign.*;
import java.foreign.memory.*;
import org.unix.curl.*;
import org.unix.curl_h;
import static org.unix.curl_h.*;
import static org.unix.easy_h.*;

public class Main {
   public static void main(String[] args) {
       try (Scope s = curl_h.scope().fork()) { 
           curl_global_init(CURL_GLOBAL_DEFAULT);
           Pointer<Void> curl = curl_easy_init();
           if(!curl.isNull()) {
               Pointer<Byte> url = s.allocateCString(args[0]);
               curl_easy_setopt(curl, CURLOPT_URL, url);
               int res = curl_easy_perform(curl);
               if (res != CURLE_OK) {
                 System.err.println("Error fetching from: " + args[0] + " ERR: " + Pointer.toString(curl_easy_strerror(res)));
                 curl_easy_cleanup(curl);
               }
           }
           curl_global_cleanup();
       }
    }
}
```

A couple of things to point out here. Notice how we cannot directly pass a Java `String` to the `curl_easy_setopt()` call. This call accepts a memory address pointer as the `url` parameter, so we first need to do a dynamic memory allocation operation using the `Scope` and pass a `Pointer` interface instance instead. As you may find in the Panama tech docs [[5]](#refs), the `Pointer` interface helps a lot when it comes to complex C-alike pointer operations like pointer arithmetic, casts, memory dereference, etc. The `Scope` manages the runtime life-cycle of dynamic allocated memory.

Alright, now armed with this knowledge, can you extend this code to write the contents of a Curl fetched web page to a file stream?

## <a name="libssh"></a>The libSSH2 App

Here's a more complete example application that utilizes *libSSH2* [[4]](#refs) to implement a simple SSH client. 

<div class="github-card" data-github="petarov/java-panama-ssh2" data-width="400" data-height="177" data-theme="default"></div>
<script src="//cdn.jsdelivr.net/github-cards/latest/widget.js"></script>

If you'd like to have a go and adjust it to run on Windows, I'll greatly appreciate a PR.

# <a name="outro"></a>Final Thoughts

A few points from my side that I learned or have been thinking about while working with Panama.

  * The `Scope` is pretty powerful. You can allocate callback pointers, structs, arrays. The concept of layouts and layout types deserves more time for me to fully explore and grasp.
  * It comes as no surprise that writing Java code using a native library is more extensive and requires extra care, especially when it comes to not forgetting to invoke cleanup API calls that the library requires.
  * I/O in most native libraries requires a file descriptor, which isn't easy to get in Java [[6]](#refs). This, however, is not directly related to the **java.foreign** API.
  * Some libraries define C++ style function prototypes without argument as opposed to C-style `Void` argument types. The Foreign Docs [[3]](#refs) have an example about this case when using the TensorFlow C API.
  * I haven't explored if it would be possible to use Panama with Go or Rust created native libraries. That would be pretty cool.

Thanks for reading!🍻

# <a name="refs"></a>References

1. [GOTO 2019 • Project Panama part](https://www.youtube.com/watch?v=vJrHHe3IbQs&t=2172)
2. [Java Native Interface](https://en.wikipedia.org/wiki/Java_Native_Interface)
3. [Panama Foreign Docs](https://hg.openjdk.java.net/panama/dev/raw-file/4810a7de75cb/doc/panama_foreign.html)
4. [libSSH2](https://www.libssh2.org) - a client-side C library implementing the SSH2 protocol
5. [Panama Binder Docs v3](https://cr.openjdk.java.net/~mcimadamore/panama/panama-binder-v3.html)
6. [Most efficient way to pass Java socket file descriptor to C binary file
](https://stackoverflow.com/q/11455803/10364676)