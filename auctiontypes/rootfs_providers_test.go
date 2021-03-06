package auctiontypes_test

import (
	"encoding/json"
	"net/url"

	"github.com/cloudfoundry-incubator/auction/auctiontypes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RootFSProviders", func() {
	var (
		arbitrary auctiontypes.ArbitraryRootFSProvider
		fixedSet  auctiontypes.FixedSetRootFSProvider
		providers auctiontypes.RootFSProviders

		providersJSON string
	)

	BeforeEach(func() {
		arbitrary = auctiontypes.ArbitraryRootFSProvider{}
		fixedSet = auctiontypes.NewFixedSetRootFSProvider("baz", "quux")
		providers = auctiontypes.RootFSProviders{
			"foo": arbitrary,
			"bar": fixedSet,
		}

		providersJSON = `{
				"foo": {
					"type": "arbitrary"
				},
				"bar": {
					"type": "fixed_set",
					"set": {"baz":{}, "quux":{}}
				}
			}`

	})

	It("serializes", func() {
		payload, err := json.Marshal(providers)
		Ω(err).ShouldNot(HaveOccurred())

		Ω(payload).Should(MatchJSON(providersJSON))
	})

	It("deserializes", func() {
		var providersResult auctiontypes.RootFSProviders
		err := json.Unmarshal([]byte(providersJSON), &providersResult)
		Ω(err).ShouldNot(HaveOccurred())

		Ω(providersResult).Should(Equal(providers))
	})

	Describe("Match", func() {
		Describe("ArbitraryRootFSProvider", func() {
			It("matches any URL", func() {
				rootFS, err := url.Parse("some://url")
				Ω(err).ShouldNot(HaveOccurred())

				Ω(arbitrary.Match(*rootFS)).Should(BeTrue())
			})
		})

		Describe("FixedSetRootFSProvider", func() {
			It("matches a URL in the set", func() {
				rootFS, err := url.Parse("some:baz")
				Ω(err).ShouldNot(HaveOccurred())

				Ω(fixedSet.Match(*rootFS)).Should(BeTrue())
			})

			It("does not match a URL not in the set", func() {
				rootFS, err := url.Parse("some://baz-not-present/here")
				Ω(err).ShouldNot(HaveOccurred())

				Ω(fixedSet.Match(*rootFS)).Should(BeFalse())
			})
		})

		Describe("RootFSProviders", func() {
			Context("for a scheme with an arbitrary provider", func() {
				It("matches any url", func() {
					rootFS, err := url.Parse("foo://any/url/is#ok")
					Ω(err).ShouldNot(HaveOccurred())

					Ω(providers.Match(*rootFS)).Should(BeTrue())
				})
			})

			Context("for a scheme with a fixed-set provider", func() {
				It("matches for a url in the set", func() {
					rootFS, err := url.Parse("bar:quux")
					Ω(err).ShouldNot(HaveOccurred())

					Ω(providers.Match(*rootFS)).Should(BeTrue())
				})

				It("does not match for a url not in the set", func() {
					rootFS, err := url.Parse("bar:quux/not?in=theset")
					Ω(err).ShouldNot(HaveOccurred())

					Ω(providers.Match(*rootFS)).Should(BeFalse())
				})
			})

			Context("for a scheme not in the map", func() {
				It("does not match", func() {
					rootFS, err := url.Parse("missingscheme://host/path")
					Ω(err).ShouldNot(HaveOccurred())

					Ω(providers.Match(*rootFS)).Should(BeFalse())
				})
			})
		})
	})
})
