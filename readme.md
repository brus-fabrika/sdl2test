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

### TODO
 [x] make pause by SPACE key - DONE
 [x] add mouse wheel to zoom in/out - DONE
 [x]  make triangle shape - was already DONE
 [ ] don't draw lines outside of the window - partially DONE - and partially fixes memory issue
 [ ] draw filled poly's with SDL_RenderGeometry or SDL_RenderGeometryRaw
 [ ] implement non-visible surfaces (lines) detection
 [ ] make text rendering more efficient
 [x] move event processing to separate function - DONE
 [ ] put engine and game to packages
