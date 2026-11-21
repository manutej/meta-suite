# ont-spec (incubating)

Specify a system before building it: a plain-language description becomes typed atoms in a
refinement tree up to four levels deep, checked for composition at every level, looped until it
converges, then mapped to typed contracts, architecture forces and an investigation brief with a
verification-first plan. See [SKILL.md](SKILL.md).

```bash
python3 scripts/ontspec.py check ontology.yaml [--converging]
python3 scripts/ontspec.py operad ontology.yaml      # needs the operadic-stitch skill installed
python3 scripts/ontspec.py delta pass-1.yaml pass-2.yaml
python3 scripts/ontspec.py contracts ontology.yaml --out contracts/
python3 scripts/ontspec.py elements ontology.yaml
python3 scripts/ontspec.py render ontology.yaml
python3 -m pytest scripts -q                          # 19 tests; 3 run the real engine when installed
```

Lives under `next/` until it has been used on more than one system; the published suite is
checked for exactly nine skills.
