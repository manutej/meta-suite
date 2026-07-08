# Internal typed model (do NOT surface to the user unless they ask to review internal work)

Model context as a persistent knowledge object K and each unit of work as a morphism that reads a
slice of K and may write back to it.

- **Substrate** fixes the representation of K (the category in which knowledge objects live);
  plain markdown is chosen because its morphisms — read, diff, grep, merge — are universal.
- **Load policy** is a functor selecting a sub-object `core ⊂ K` always supplied, with references
  pulled by a task-indexed family of inclusions (progressive disclosure = a colimit assembled only
  as needed).
- **Compression** is an idempotent endofunctor on K that quotients duplicates and discards stale
  objects, keeping |K| proportional to distinct knowledge rather than activity.
- **Persistence/handoff** is the write-back morphism closing each work loop into K, so the next
  morphism starts from an enriched object (warm start).
- **Discoverability** is the natural transformation guaranteeing every object of K is reachable from
  the index object; without it, written knowledge is unreachable and does not compound.

This framing is a private reasoning aid. Never present it, name the structure, or ask the user about
it unless they explicitly request the internal work. The user sees a plain context plan and a
droppable-in scaffold.
