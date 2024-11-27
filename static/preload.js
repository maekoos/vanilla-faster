window.addEventListener('load', () => {
    const cachedUrls = new Set();

    const els = document.querySelectorAll("a[href]")
    for (const el of els) {
        const href = el.getAttribute('href');
        const url = new URL(href, window.location);
        if (url.origin !== window.location.origin) continue;
        if (url.href.slice(url.origin.length) === window.location.href.slice(window.location.origin.length)) continue;

        let isOver = false;
        const handleOver = async (ev) => {
            if (cachedUrls.has(url.href)) return;
            isOver = true;

            console.log(`Prefetching ${url.href}...`);
            const dom = await fetch(url.href)
                .then(x => x.text())
                .then(x => new DOMParser().parseFromString(x, "text/html"));

            const subEls = dom.querySelectorAll('img[src]');
            for (const el of subEls) {
                // addPrefetch(el.src);
                // Await a manual prefetch to throttle a bit
                await manualPrefetch(el.src);
                if (isOver == false) break;
            }
        };

        // Stop loading when mouse leaves
        const handleLeave = async (ev) => {
            isOver = false;
        };

        el.addEventListener('mouseenter', handleOver);
        el.addEventListener('mouseleave', handleLeave);
    }

    for (const lazyEl of document.querySelectorAll('img[loading=lazy]')) {
        // addPrefetch(lazyEl.src);
        // manualPrefetch(lazyEl.src);
        lazyEl.loading = 'eager';
    }


    function addPrefetch(src) {
        if (cachedUrls.has(src)) return;
        console.log(`Prefetching image ${src}...`);
        const linkEl = document.createElement('link');
        linkEl.rel = 'prefetch';
        linkEl.href = src;
        document.body.append(linkEl);
        cachedUrls.add(src);
    }

    async function manualPrefetch(src) {
        if (cachedUrls.has(src)) return;
        console.log(`Prefetching image ${src}...`);
        cachedUrls.add(src);
        await fetch(src);
    }
});
