### How to run the project
- Clone the project
- Install 3rd party dependencies:
    - sudo apt install libsdl2-dev libsdl2-ttf-dev libsdl2-gfx-dev
- Run in the root of the project: go mod tidy



### Questions
- Why is memory consumption ~5GB when changed Zbase (but not zoom f) ????
- Why does performance drop when zoomed in? - changing perspective (Zbase) leads to memory consumption (see prev point)
- Why does shape not rotate around center point?
- How NOT to clone mesh on every frame?

### TODO General
[x] make pause by SPACE key
[x] add mouse wheel to zoom in/out
[x] move event processing to separate function
[ ] put engine and game to packages
[ ] make build root and makefiles to build outside of sources

### TODO Rendering
[x] implement non-visible surfaces (lines) detection
[x] make triangle shape
[x] draw filled poly's with SDL_RenderGeometry or SDL_RenderGeometryRaw
[ ] don't draw lines outside of the window - partially DONE - and partially fixes memory issue
[ ] make text rendering more efficient
[ ] make hidden faces behind invisible (not showing, Z-order)
[ ] implement simple lighting model for faces
[ ] implement camera shift/move
