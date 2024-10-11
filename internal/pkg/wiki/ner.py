import sys
import spacy
import json

nlp = spacy.load("en_core_web_sm")

ENTITY_LABELS = {
    "ORG": "Developer",   # Organizations -> Developers
    "PRODUCT": "Platform",  # Products -> Platforms
    "NORP": "Genre",       # Nationalities, religious or political groups -> Genres (as in game genres)
}

def extract_entities(text):
    doc = nlp(text)
    entities = []

    for ent in doc.ents:
        label = ENTITY_LABELS.get(ent.label_, "Other")
        entities.append({
            "text": ent.text,
            "label": label
        })
    return entities

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python ner.py 'Your text here'")
        sys.exit(1)

    input_text = sys.argv[1]
    entities = extract_entities(input_text)

    # Output the entities as a JSON string
    print(json.dumps(entities))
