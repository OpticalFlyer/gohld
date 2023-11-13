function setVH() {
    let vh = window.innerHeight * 0.01;
    document.documentElement.style.setProperty('--vh', `${vh}px`);
    }

    window.addEventListener('resize', setVH);
    window.addEventListener('orientationchange', setVH);

    // Initial set
    setVH();