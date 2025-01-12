### Questions
 - Why is memory consumption ~5GB when changed Zbase (but not zoom f) ????
 - Why does performance drop when zoomed in? - changing perspective (Zbase) leads to memory consumption (see prev point)
 - Why does shape not rotate around center point?
 - How NOT to clone mesh on every frame?

### TODO
 - make pause by SPACE key - DONE
 - add mouse wheel to zoom in/out - DONE
 - make triangle shape - was already DONE
 - don't draw lines outside of the window - partially DONE - and partially fixes memory issue
 - draw filled poly's with SDL_RenderGeometry or SDL_RenderGeometryRaw
 - implement non-visible lines detection
 - make text rendering more efficient
 - move event processing to separate function
