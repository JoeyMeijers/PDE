import json
from collections import defaultdict


def run(input_data: str) -> str:
    """
    input_data: JSON string met lijst van dicts, bijv:
    [{"age": 25, "score": 88}, ...]

    output: JSON string met gemiddelde score per leeftijdsgroep
    [{"age_group": 20, "avg_score": 75.5}, ...]
    """
    data = json.loads(input_data)

    groups = defaultdict(list)
    for row in data:
        age = int(row["age"])
        score = float(row["score"])
        group = (age // 10) * 10
        groups[group].append(score)

    result = []
    for group, scores in sorted(groups.items()):
        avg = sum(scores) / len(scores)
        result.append({"age_group": group, "avg_score": round(avg, 2)})

    return json.dumps(result)
