package spinner

import (
	"os"
	"strings"
	"unicode"
	"fmt"
	"golang.org/x/text/unicode/rangetable"
	
	"github.com/diggernaut/whatlanggo"
	"github.com/diggernaut/wnram"
	"github.com/jdkato/prose/tag"
	//	"github.com/kljensen/snowball"
)

var (
	wn *wnram.Handle
	err error
	pos map[string]wnram.PartOfSpeech
)

func init() {
	wn, err = wnram.New("E:/Jobs/Current/Smartologic/Wordnet")
	if err != nil {
		fmt.Printf("Cannot init Wordnet: %v\n", err)
		os.Exit(1)
	}
	pos = make(map[string]wnram.PartOfSpeech)
	// pos["CC"] = nil // Coordinating conjunction
	// pos["CD"] = nil // Cardinal number
	// pos["DT"] = nil // Determiner
	// pos["EX"] = nil // Existential
	// pos["FW"] = nil // Foreign word
	// pos["IN"] = nil // Preposition or subordinating conjunction
	pos["JJ"] = wnram.Adjective
	pos["JJR"] = wnram.Adjective
	pos["JJS"] = wnram.Adjective
	// pos["LS"] = nil // List item marker
	// pos["MD"] = nil // Modal
	pos["NN"] = wnram.Noun
	pos["NNS"] = wnram.Noun
	pos["NNP"] = wnram.Noun
	pos["NNPS"] = wnram.Noun
	// pos["PDT"] = nil // Predeterminer
	// pos["POS"] = nil // Possessive ending
	// pos["PRP"] = nil // Personal pronoun
	// pos["PRP$"] = nil // Possessive pronoun
	pos["RB"] = wnram.Adverb
	pos["RBR"] = wnram.Adverb
	pos["RBS"] = wnram.Adverb
	// pos["RP"] = nil // Particle
	// pos["SYM"] = nil // Symbol
	// pos["TO"] = nil // to
	// pos["UH"] = nil // Interjection
	pos["VB"] = wnram.Verb // base form
	pos["VBD"] = wnram.Verb // past tense
	pos["VBG"] = wnram.Verb // gerund or present participle
	pos["VBN"] = wnram.Verb // past participle
	pos["VBP"] = wnram.Verb // non-3rd person singular present
	pos["VBZ"] = wnram.Verb // 3rd person singular present
	// pos["WDT"] = nil // Wh-determiner
	// pos["WP"] = nil // Wh-pronoun
	// pos["WP$"] = nil // Possessive wh-pronoun
	// pos["WRB"] = nil // Wh-adverb
}
// Spintax makes spitax template from given article
func Spintax(article string) string {
	blocks := splitByBlocks(article)
	lang := whatlanggo.Langs[whatlanggo.DetectLang(article)]

	for _, block := range blocks {
		tokens := tokenize(block["text"], lang)
		fmt.Printf("S: %v\n", tokens)
	}
	return "OK"
}

func splitByBlocks(s string) []map[string]string {
	var blocks []map[string]string
	var block string
	allowedChars := rangetable.New(rune('\''))

	for _, letter := range s {
		if unicode.IsLetter(letter) || unicode.IsSpace(letter) || unicode.Is(allowedChars,letter) {
			if !(unicode.IsSpace(letter) && block == "") {
				block += string(letter)
			}
		} else {
			bl := make(map[string]string)
			bl["text"] = block
			bl["end"] = string(letter)
			blocks = append(blocks, bl)
			block = ""
		}
	}

	return blocks
}

func tokenize(s string, lang string) []map[string]string {
	var tokens []map[string]string
	
	tagger := tag.NewPerceptronTagger()
	tags := tagger.Tag(strings.Fields(s))
	syns := make(map[string][]string)
	graphs := make(map[string][]string)
	layers := [][]wnram.ID{}
	for _, tag := range tags {
		fmt.Println(tag.Text, tag.Tag)
		if p, ok := pos[tag.Tag]; ok {
			criteria := wnram.Criteria{Matching: tag.Text, POS: []wnram.PartOfSpeech{p}}
			synsets, err := wn.Lookup(criteria)
			layer := []wnram.ID{}
			if err == nil {
				for _, synset := range synsets {
					nodeID := synset.NodeID()
					layer = append(layer, nodeID)
					synonyms := synset.Synonyms()
					syns[nodeID.String()] = synonyms
					if graphs[tag.Text] == nil {
						graphs[tag.Text] = []string{}
					}
					graphs[tag.Text] = append(graphs[tag.Text], nodeID.String())
				}
			} else {
				fmt.Printf("err: %v\n", err)
			}
			layers = append(layers, layer)
		}
	}
	for i := 0; i < len(layers); i++ {
		for j := i + 1; j < len(layers); j++ {
			for _, nodeID1 := range layers[i] {
				for _, nodeID2 := range layers[j] {
					weight := wn.GetDistance(nodeID1, nodeID2)
					fmt.Printf("Weight: %v\n", weight)
				}
			}
			fmt.Printf("Iterj: %v\n", j)
		}
		fmt.Printf("Iteri: %v\n", i)
	}
	os.Exit(1)
	
	
	return tokens
}


func splitOnNonLetters(s string) []string {
	allowedChars := rangetable.New(rune('\''))
	notALetter := func(char rune) bool { return !(unicode.IsLetter(char) || unicode.IsSpace(char) || unicode.Is(allowedChars,char)) }
	return strings.FieldsFunc(s, notALetter)
}

/*
def getSpintax(self, text):
sentences = self.splitToSentences(text)
stemmer = PorterStemmer()
spintax = ""
for sentence in sentences:
	tokens = regexp_tokenize(sentence, "[\w']+")
	for token in tokens:
		stem = stemmer.stem(token)
		n, syn = self.getSynonyms(stem)
		spintax += "{"
		spintax += token
		spintax += "|"
		for x in range(n):
			spintax += syn[x]
			if x < n-1:
				spintax += "|"
			else:
				spintax += "} "
return spintax
*/