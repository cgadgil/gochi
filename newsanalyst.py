from newsapi import NewsApiClient
import nltk
from nltk.tokenize import wordpunct_tokenize
from nltk.probability import FreqDist
from nltk.stem import WordNetLemmatizer 
from nltk.corpus import stopwords
from nltk.stem import PorterStemmer
from nltk.tokenize import RegexpTokenizer
import sys

def clean_text(text):
    "Clean up text - lemmatize/stem/remove stop words"
    if not text:
        return []
    stop_words = set(stopwords.words('english'))
    #print("Text: " + text)
    word_tokens = wordpunct_tokenize(text) 
    filtered_sentence = []
    for w in word_tokens:
        if w not in stop_words:
            if w.isalpha() and w != "chars":
                filtered_sentence.append(w)
    return filtered_sentence


# Init
newsapi = NewsApiClient(api_key='0c2182c4b4d040d6824e5684344e06a6')

# /v2/top-headlines
et = newsapi.get_everything(q=sys.argv[1], language='en')

fdist = FreqDist()

for article in et['articles']:
    import pdb
    #pdb.set_trace()
    cleaned_tokens = clean_text(article["content"])
    for word in cleaned_tokens:
        fdist[word.lower()] += 1

print(fdist.most_common(50))
