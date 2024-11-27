# VanillaFaster

It is silly how much attention [McMaster-Carr](https://www.mcmaster.com/) and [NextFaster](https://github.com/ethanniser/NextFaster) has gotten recently. Yes, they are fast. But not *that* fast. Sufficient cache policies, and a little [script](./static/preload.js) for preloading images and links on hover is all it takes to get a pretty fast e-commerce website - just look at this super simple website.

Yes, it is written in a fast language, no that is not what makes it fast. If you are interested in what makes this website fast, take a look at the caching, the [resizing of images](./resize/resize.go) and the [preload script](./static/preload.js) - that's it.

## Things could be better
This is just a quick and dirty showcase of how easy it is to make a fast (simple e-commerce) website - there are a lot of improvements that can be made. Here are a few things I haven't done:
- Error handling
- Caching database results
- Resize may be a great subject for DDOS (doesn't limit the url nor the maxSize property - thus a great way to (d)dos the server is to use `http://.../_image?url=${encode_url("https://localhost:3000/the-largest-legit-image.jpg?"+random_string())}&size=${randInt(100,200)}`)
- Pre-fetching images and storing them locally (only cached in memory for now)
- Searching, auth and basket
- Design...
- etc.


## One note

I'd like to note that performance isn't everything - and often comes at a cost. In this case, a lot of links and images are prefetched resulting in high bandwidth usage on both client and server and generally higher load on the server.


## Running
*This website uses the exact same database as [NextFaster](https://github.com/ethanniser/NextFaster) - set it up and come back.*


```
./db.sh
go generate && go run .
```