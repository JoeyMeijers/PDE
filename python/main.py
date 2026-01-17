import sys
import json

data = json.load(sys.stdin)

data["cleaned"] = True

json.dump(data, sys.stdout)
