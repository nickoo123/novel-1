(function(){
    var bp = document.createElement('script');
    var curProtocol = window.location.protocol.split(':')[0];
    if (curProtocol === 'https') {
        bp.src = 'https://zz.bdstatic.com/linksubmit/push.js';
    }
    else {
        bp.src = 'http://push.zhanzhang.baidu.com/push.js';
    }
    var s = document.getElementsByTagName("script")[0];
    s.parentNode.insertBefore(bp, s);
})();
window.Echo=(function(b,h,d){var i=[];var k=function(){};var f,j,g;var c=function(n){var o=n.getBoundingClientRect();return((o.top>=0&&o.left>=0&&o.top)<=(window.innerHeight||h.documentElement.clientHeight)+f)};var e=function(){var q=i.length;if(q>0){for(var p=0;p<q;p++){var n=i[p];if(n&&c(n)){var o=n.src;n.onerror=function(){this.onerror=null;this.src=o};var r=n.getAttribute("data-echo");if(r){n.src=r;if(n.naturalHeight&&n.naturalHeight===0){n.src=o}}n.removeAttribute("data-echo");k(n);i.splice(p,1);q=i.length;p--}}}else{if(h.removeEventListener){b.removeEventListener("scroll",m)}else{b.detachEvent("onscroll",m)}clearTimeout(g)}};var m=function(){clearTimeout(g);g=setTimeout(e,j)};var l=function(q){var n=h.querySelectorAll("img[data-echo]");var p=q||{};f=parseInt(p.offset||0);j=parseInt(p.throttle||250);k=p.callback||k;for(var o=0;o<n.length;o++){i.push(n[o])}e();if(h.addEventListener){b.addEventListener("scroll",m,false);b.addEventListener("load",m,false);b.addEventListener("touchmove",m,false)}else{b.attachEvent("onscroll",m);b.attachEvent("onload",m);b.attachEvent("ontouchmove",m)}};var a=function(){if(h.removeEventListener){b.removeEventListener("scroll",m)}else{b.detachEvent("onscroll",m)}clearTimeout(g);i=[]};return{init:l,destroy:a,render:e}})(this,document);