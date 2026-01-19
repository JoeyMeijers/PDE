PDP — Polyglot Data Pipeline Engine

ISE is een taal-agnostische, container-gebaseerde data pipeline engine. Het voert gestructureerde data transformaties uit via een flexibele pipeline van functies, ongeacht of deze in Python, Rust of andere talen zijn geschreven.

Belangrijkste kenmerken:
	•	Pipeline-gebaseerd: definieer een reeks stappen die data transformeren, filteren of verrijken.
	•	Polyglot / Containerized: elke stap kan in een andere programmeertaal worden uitgevoerd via Docker, volledig geïsoleerd.
	•	JSON-georiënteerd: data wordt tussen stappen doorgegeven als JSON, wat compatibiliteit en uitbreidbaarheid vergroot.
	•	Extensibel: eenvoudig nieuwe functies toevoegen in Python, Rust of andere talen zonder de core engine te wijzigen.
	•	Logging & observability: elke stap wordt gelogd met starttijd, duur, output size en fouten.

Gebruiksscenario’s:
	•	ETL-processen voor grote datasets
	•	Multi-language function orchestratie
	•	Prototyping van datatransformaties en statistische berekeningen
	•	Stress- en performance-tests van pipelines
