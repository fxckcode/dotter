// Package tlds provides a curated list of TLDs and parallel scanning.
package tlds

// DefaultTLDs returns the default curated list of TLDs to check.
// These are the most popular and commonly used extensions.
func DefaultTLDs() []string {
	return []string{
		".com", ".io", ".dev", ".app", ".tech", ".ai",
		".net", ".org", ".co", ".me", ".xyz", ".info",
		".design", ".tools", ".cloud", ".digital", ".network",
		".world", ".online", ".site", ".store", ".pro",
		".blog", ".agency", ".media", ".email", ".social",
		".studio", ".works", ".today", ".live", ".run",
		".systems", ".software", ".engineering", ".codes",
		".wiki", ".guide", ".report", ".finance",
		".legal", ".law", ".health", ".life",
		".team", ".group", ".company", ".enterprises",
		".uno", ".win", ".bid", ".trade", ".webcam",
		".review", ".science", ".download", ".party",
		".date", ".racing", ".accountant",
	}
}

// ExtendedTLDs returns a longer list including all default + niche TLDs.
func ExtendedTLDs() []string {
	extra := []string{
		".pics", ".pictures", ".photo", ".photos", ".image",
		".gallery", ".art", ".graphics", ".video", ".audio",
		".music", ".fans", ".mom", ".lol", ".fun",
		".games", ".gaming", ".play", ".zone", ".space",
		".rocks", ".news", ".press", ".report", ".review",
		".marketing", ".adult", ".porn", ".sex",
		".xxx", ". dating", ".faith", ".church",
		".bible", ".islam", ".hindu", ".jewelry",
		".gold", ".diamonds", ".gift", ".gifts",
		".flowers", ".toys", ".tires", ".glass",
		".cafe", ".pizza", ".coffee", ".beer", ".wine",
		".vodka", ".whisky", ".fish", ".seafood",
		".restaurant", ".menu", ".recipes", ".diet",
		".fitness", ".guru", ".expert", ".coach",
		".training", ".school", ".college", ".university",
		".academy", ".courses", ".education", ".degree",
		".doctor", ".dentist", ".vet", ".hospital",
		".clinic", ".surgery", ".pharmacy", ".medical",
		".travel", ".vacations", ".holiday", ".tours",
		".cruises", ".flights", ".hotels", ".rentals",
		".camp", ".beach", ".surf", ".ski",
		".horse", ".rodeo", ".farm", ".garden",
		".land", ".house", ".build", ".construction",
		".archi", ".engineer", ".contractors",
	}
	return append(DefaultTLDs(), extra...)
}
