function animateText() {
  var ascii = document.getElementById("ascii");
  var text = ascii.textContent;
  var to = text.length,
    from = 0;

  animate({
    duration: 1000,
    timing: bounce,
    draw: function (progress) {
      var result = (to - from) * progress + from;
      ascii.textContent = text.slice(0, Math.ceil(result));
    },
  });
}

function bounce(x) {
  // return 1 - Math.cos((x * Math.PI) / 2);
  return x;
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

function bufferDecode(b64) {
  const pad = "=".repeat((4 - (b64.length % 4)) % 4);
  const base64 = (b64 + pad).replace(/-/g, "+").replace(/_/g, "/");
  const raw = atob(base64);
  return Uint8Array.from([...raw].map((c) => c.charCodeAt(0)));
}

function bufferEncode(buffer) {
  const bytes = new Uint8Array(buffer);
  const bin = String.fromCharCode(...bytes);
  const base64 = btoa(bin);
  return base64.replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");
}

async function setupCredentials() {
  const { publicKey } = await fetch("/sudo/registration/begin").then((r) =>
    r.json(),
  );

  publicKey.challenge = bufferDecode(publicKey.challenge);
  publicKey.user.id = bufferDecode(publicKey.user.id);
  if (publicKey.excludeCredentials) {
    for (var i = 0; i < publicKey.excludeCredentials.length; i++) {
      publicKey.excludeCredentials[i].id = bufferDecode(
        publicKey.excludeCredentials[i].id,
      );
    }
  }
  console.log(publicKey);

  var cred;
  try {
    cred = await navigator.credentials.create({ publicKey });
    console.log(cred);
  } catch (err) {
    if (err.name === "InvalidStateError") {
      alert("Устройство уже зарегистрировано");
      return;
    }
    console.error(err);
  }

  const response2 = await fetch("/sudo/registration/finish", {
    method: "POST",
    body: JSON.stringify({
      id: cred.id,
      rawId: bufferEncode(cred.rawId),
      type: cred.type,
      response: {
        clientDataJSON: bufferEncode(cred.response.clientDataJSON),
        attestationObject: bufferEncode(cred.response.attestationObject),
      },
    }),
  });

  if (response2.redirected) {
    window.location.href = response2.url;
  } else {
    console.error(response2);
  }
}

async function loginCredentials() {
  const { publicKey } = await fetch("/sudo/login/begin").then((r) => r.json());

  publicKey.challenge = bufferDecode(publicKey.challenge);
  publicKey.allowCredentials.forEach((item) => {
    item.id = bufferDecode(item.id);
  });
  console.log(publicKey);

  var cred;
  try {
    cred = await navigator.credentials.get({ publicKey });
    console.log(cred);
  } catch (err) {
    if (err.name === "InvalidStateError") {
      alert("Устройство уже зарегистрировано");
      return;
    }
    console.error(err);
  }

  const response2 = await fetch("/sudo/login/finish", {
    method: "POST",
    body: JSON.stringify({
      id: cred.id,
      rawId: bufferEncode(cred.rawId),
      type: cred.type,
      response: {
        authenticatorData: bufferEncode(cred.response.authenticatorData),
        clientDataJSON: bufferEncode(cred.response.clientDataJSON),
        signature: bufferEncode(cred.response.signature),
        userHandle: bufferEncode(cred.response.userHandle),
      },
    }),
    credentials: "include",
  });
  if (response2.redirected) {
    window.location.href = response2.url;
  } else {
    console.error(response2);
  }
}
