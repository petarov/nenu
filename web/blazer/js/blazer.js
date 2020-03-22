;(function(w) {
  function gs(k, def) {
    if (localStorage) {
      return localStorage.getItem(k) || def;
    }
    return def;
  }
  function ss(k, v) {
    if (localStorage) {
      localStorage.setItem(k, v);
    }
  }
  var bumpIt = function(el, mode) {
    if (mode=='blazer') {
      el.innerText = 'light mode';
    } else {
      el.innerText = 'dark mode';
    }
    var link = document.getElementById("theme-css");
    link.href = 'styles/sakura-' + mode + '.css';
    ss('mode', mode);
  };
  var bumper = document.getElementById("cuk");
  bumper.addEventListener('click', function(e) {
    bumpIt(e.target, gs('mode', 'blazer') =='blazer' ? 'light' : 'blazer');
  });
  bumpIt(bumper, gs('mode', 'blazer'));
})(window);