package common

import (
	"encoding/xml"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type SitemapGenerator struct {
	repository *Repository
	config *SitemapGeneratorConfig
}

func NewSitemapGenerator(repository *Repository) *SitemapGenerator {
	return &SitemapGenerator{
		repository: repository,
		config: &config.SitemapGeneratorConfig,
	}
}

func (this *SitemapGenerator) generateLink(codInmatriculare string, numeFirma string) string {
	codInmatriculare = strings.ReplaceAll(codInmatriculare, "/", "-")

	u := url.URL{
		Scheme: "https",
		Host:   "openfirme.ro",
		Path: "profile/" + codInmatriculare,
	}

	q := u.Query()
	q.Set("nume_firma", numeFirma)
	u.RawQuery = q.Encode()

	return u.String()
}

type XMLUrl struct {
	Loc string `xml:"loc"`
	LastMod string `xml:"lastmod"`
}

type XMLUrlSitemap struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS string   `xml:"xmlns,attr"`
	Urls []XMLUrl `xml:"url"`
}

func (this *SitemapGenerator) generateSitemapsFirme(firme []CodInmatriculareNumeFirma) {
	totalCountFirme := len(firme)

	batchSize := totalCountFirme / this.config.SitemapFirmeCount
	if totalCountFirme % this.config.SitemapFirmeCount != 0 {
		batchSize += 1
	}

	sitemapNameBase := filepath.Join(this.config.OutputDirectory, "sitemapfirme")

	lastMod := time.Now().Format("2006-01-02")

	var sitemap XMLUrlSitemap

	for index, firma := range firme {
		sitemap.Urls = append(sitemap.Urls, XMLUrl {
			Loc: this.generateLink(firma.codInmatriculare, firma.numeFirma),
			LastMod: lastMod,
		})

		if (index + 1) % batchSize == 0 {
			serialized, err := xml.Marshal(sitemap)
			if err != nil {
				panic(err)
			}

			err = os.WriteFile(sitemapNameBase + strconv.Itoa(index / batchSize) + ".xml", serialized, 0644)
			if err != nil {
				panic(err)
			}

			sitemap = XMLUrlSitemap {
				XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
				Urls: []XMLUrl{},
			}
		}
	}

	// some of them remained
	if len(sitemap.Urls) != 0 {
		serialized, err := xml.Marshal(sitemap)
		if err != nil {
			panic(err)
		}

		err = os.WriteFile(sitemapNameBase + strconv.Itoa((totalCountFirme - 1) / batchSize) + ".xml", serialized, 0644)
		if err != nil {
			panic(err)
		}
	}
}

type XMLSitemap struct {
	Loc string `xml:"loc"`
	LastMod string `xml:"lastmod"`
}

type XMLIndexSitemap struct {
	XMLName xml.Name `xml:"sitemapindex"`
	XMLNS string   `xml:"xmlns,attr"`
	Sitemaps []XMLSitemap `xml:"sitemap"`
}

func (this *SitemapGenerator) generateSitemapIndex(firme []CodInmatriculareNumeFirma) {
	lastMod := time.Now().Format("2006-01-02")

	sitemapIndex := XMLIndexSitemap {
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Sitemaps: []XMLSitemap {
			XMLSitemap {
				Loc: "https://openfirme.ro/public/assets/sitemaptop.xml",
				LastMod: lastMod,
			},
		},
	}

	for i := 0; i < this.config.SitemapFirmeCount; i++ {
		sitemapIndex.Sitemaps = append(sitemapIndex.Sitemaps, XMLSitemap {
			Loc: "https://openfirme.ro/public/assets/sitemapfirme" + strconv.Itoa(i) + ".xml",
			LastMod: lastMod,
		})
	}

	serialized, err := xml.MarshalIndent(sitemapIndex, "", "  ")
	if err != nil {
		panic(err)
	}

	outputPath := filepath.Join(this.config.OutputDirectory, "sitemap.xml")
	err = os.WriteFile(outputPath, serialized, 0644)
	if err != nil {
		panic(err)
	}
}

func (this *SitemapGenerator) GenerateSitemap() {
	if !this.config.Enabled {
		return
	}

	os.Mkdir(this.config.OutputDirectory, 0755)

	firme := this.repository.GetAllFirmeActive()

	this.generateSitemapsFirme(firme)
	this.generateSitemapIndex(firme)
}
