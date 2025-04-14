// dom content loaded
document.addEventListener("DOMContentLoaded", () => {
  function animateText() {
    var ascii = document.getElementById("ascii");
    var text = ascii.textContent;
    var to = text.length,
      from = 0;

    animate({
      duration: 5000,
      timing: bounce,
      draw: function (progress) {
        var result = (to - from) * progress + from;
        ascii.textContent = text.slice(0, Math.ceil(result));
      },
    });
  }

  function bounce(timeFraction) {
    return timeFraction;
  }

  function animate({ duration, draw, timing }) {
    var start = performance.now();

    const raf = requestAnimationFrame(function animate(time) {
      var timeFraction = (time - start) / duration;
      if (timeFraction > 1) {
        timeFraction = 1;
        cancelAnimationFrame(raf);
      }

      var progress = timing(timeFraction);

      draw(progress);

      if (timeFraction < 1) {
        requestAnimationFrame(animate);
      }
    });
  }

  // animateText();
});
