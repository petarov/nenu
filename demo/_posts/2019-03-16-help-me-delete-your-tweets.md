---
title: Help Me Delete Your Tweets!
layout: post
date: 2019-03-16T12:00:00.00Z
publish: true
---
Ok, I do apologize for the clickbaity title and to make up for it here's what this is all about in brief.


**tl;dr**, I'm looking for hackers to help me improve a Python script.


A couple of years ago I wanted to have all my tweets and likes deleted, but still keep them somewhere for private purposes, mainly for history reference and personal archive. Because I'm not really that active on Twitter, deleting my account would probably have made more sense. But hey! I still wanted to put my uninformed opinion out there from time to time.


Long story short, I hacked this small Python command line tool that consumes the Twitter API. It finds all tweets and/or likes after a specified date-time point, puts them in an ePub book and then (optionally) removes them from Twitter. Because it's just a command line thing, it's also simple to have it run in a cron job as well. So far so good, it seems to work alright, however, there're some limitations. I couldn't figure out how to save tweet images with [Tweepy](http://www.tweepy.org/) and the ePub E-book generation could use more options like front and back cover images, table of contents, layout options, etc.


So I guess this is a call for contributions. If you're looking for some open source project to contribute to and this looks fun to you, you're more than welcome to join in. I'm actually not a Python developer, but a Python hacker, so this is more like a *hack-it-for-fun-but-still-make-it-useful* thing here 


Here's the GitHub repo and thanks for reading!

<div class="github-card" data-github="petarov/shut-up-bird" data-width="400" data-height="177" data-theme="default"></div>
<script src="//cdn.jsdelivr.net/github-cards/latest/widget.js"></script>