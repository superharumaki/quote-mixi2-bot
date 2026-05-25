import csv
import json

with open('quotes.csv', encoding='utf-16') as f:
    reader = csv.DictReader(f, delimiter='\t')
    data = list(reader)

result = []

for row in data:
    result.append({
        'text': row['text'],
        'author': row['author']
    })

with open('quotes.json', 'w', encoding='utf-8') as f:
    json.dump(result, f, ensure_ascii=False, indent=2)

print("quotes.json を更新しました")
