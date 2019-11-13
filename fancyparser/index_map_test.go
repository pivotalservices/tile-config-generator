package fancyparser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/tile-config-generator/fancyparser"
)

var _ = Describe("IndexMap", func() {

	Context("Filtering keys", func() {
		var indexMap IndexMap

		Context("IsPlaceholder", func() {
			Context("When a placeholder is provided", func() {
				It("returns true", func() {
					variable := "((beep-boop))"
					Expect(IsPlaceholder(variable)).To(BeTrue())
				})
			})

			Context("When a placeholder is not provided", func() {
				It("returns false", func() {
					variable := "beep-boop))"
					Expect(IsPlaceholder(variable)).To(BeFalse())

					variable = "((beep-boop"
					Expect(IsPlaceholder(variable)).To(BeFalse())

					variable = "be((ep-b))oop"
					Expect(IsPlaceholder(variable)).To(BeFalse())
				})
			})
		})

		BeforeEach(func() {
			indexMap = IndexMap{
				"((networking_poe_ssl_certs_0/certificate))": []Index{
					Index{Type: IndexTypeMap, MapIndex: ".properties.networking_poe_ssl_certs"},
					Index{Type: IndexTypeMap, MapIndex: "value"},
					Index{Type: IndexTypeList, ListIndex: 0},
					Index{Type: IndexTypeMap, MapIndex: "certificate"},
					Index{Type: IndexTypeMap, MapIndex: "cert_pem"},
				},
				"((networking_poe_ssl_certs_0/privatekey))": []Index{
					Index{Type: IndexTypeMap, MapIndex: ".properties.networking_poe_ssl_certs"},
					Index{Type: IndexTypeMap, MapIndex: "value"},
					Index{Type: IndexTypeList, ListIndex: 0},
					Index{Type: IndexTypeMap, MapIndex: "certificate"},
					Index{Type: IndexTypeMap, MapIndex: "private_key_pem"},
				},
				"BOB": []Index{
					Index{Type: IndexTypeMap, MapIndex: ".properties.networking_poe_ssl_certs"},
					Index{Type: IndexTypeMap, MapIndex: "value"},
					Index{Type: IndexTypeList, ListIndex: 0},
					Index{Type: IndexTypeMap, MapIndex: "name"},
				},
				"SAM": []Index{
					Index{Type: IndexTypeMap, MapIndex: ".properties.networking_poe_ssl_certs"},
					Index{Type: IndexTypeMap, MapIndex: "value"},
					Index{Type: IndexTypeList, ListIndex: 0},
					Index{Type: IndexTypeMap, MapIndex: "name"},
				},
			}
		})

		Context("GetPlaceholderValueIndexes", func() {
			It("Filters out only indexes with placeholder keys and strips parenths", func() {
				Expect(indexMap.GetPlaceholderValueIndexes()).To(Equal(IndexMap{
					"networking_poe_ssl_certs_0/certificate": []Index{
						Index{Type: IndexTypeMap, MapIndex: ".properties.networking_poe_ssl_certs"},
						Index{Type: IndexTypeMap, MapIndex: "value"},
						Index{Type: IndexTypeList, ListIndex: 0},
						Index{Type: IndexTypeMap, MapIndex: "certificate"},
						Index{Type: IndexTypeMap, MapIndex: "cert_pem"},
					},
					"networking_poe_ssl_certs_0/privatekey": []Index{
						Index{Type: IndexTypeMap, MapIndex: ".properties.networking_poe_ssl_certs"},
						Index{Type: IndexTypeMap, MapIndex: "value"},
						Index{Type: IndexTypeList, ListIndex: 0},
						Index{Type: IndexTypeMap, MapIndex: "certificate"},
						Index{Type: IndexTypeMap, MapIndex: "private_key_pem"},
					},
				}))
			})
		})

		Context("FilterHardcodedValuesIndexes", func() {
			It("Filters out only indexes with hardcoded value keys", func() {
				Expect(indexMap.GetHardcodedValueIndexes()).To(Equal(IndexMap{
					"BOB": []Index{
						Index{Type: IndexTypeMap, MapIndex: ".properties.networking_poe_ssl_certs"},
						Index{Type: IndexTypeMap, MapIndex: "value"},
						Index{Type: IndexTypeList, ListIndex: 0},
						Index{Type: IndexTypeMap, MapIndex: "name"},
					},
					"SAM": []Index{
						Index{Type: IndexTypeMap, MapIndex: ".properties.networking_poe_ssl_certs"},
						Index{Type: IndexTypeMap, MapIndex: "value"},
						Index{Type: IndexTypeList, ListIndex: 0},
						Index{Type: IndexTypeMap, MapIndex: "name"},
					},
				}))
			})
		})
	})

	Context("GetMapWithPrependedIndex", func() {
		var indexMap IndexMap

		BeforeEach(func() {
			indexMap = IndexMap{
				"((networking_poe_ssl_certs_0/certificate))": []Index{
					Index{Type: IndexTypeMap, MapIndex: "value"},
					Index{Type: IndexTypeList, ListIndex: 0},
					Index{Type: IndexTypeMap, MapIndex: "certificate"},
					Index{Type: IndexTypeMap, MapIndex: "cert_pem"},
				},
				"((networking_poe_ssl_certs_0/privatekey))": []Index{
					Index{Type: IndexTypeMap, MapIndex: "value"},
					Index{Type: IndexTypeList, ListIndex: 0},
					Index{Type: IndexTypeMap, MapIndex: "certificate"},
					Index{Type: IndexTypeMap, MapIndex: "private_key_pem"},
				},
			}
		})

		It("returns a map with the provided index prepended to all the indexes", func() {
			newMap := indexMap.GetMapWithPrependedIndex(Index{Type: IndexTypeMap, MapIndex: ".properties.networking_poe_ssl_certs"})
			Expect(newMap).To(Equal(IndexMap{
				"((networking_poe_ssl_certs_0/certificate))": []Index{
					Index{Type: IndexTypeMap, MapIndex: ".properties.networking_poe_ssl_certs"},
					Index{Type: IndexTypeMap, MapIndex: "value"},
					Index{Type: IndexTypeList, ListIndex: 0},
					Index{Type: IndexTypeMap, MapIndex: "certificate"},
					Index{Type: IndexTypeMap, MapIndex: "cert_pem"},
				},
				"((networking_poe_ssl_certs_0/privatekey))": []Index{
					Index{Type: IndexTypeMap, MapIndex: ".properties.networking_poe_ssl_certs"},
					Index{Type: IndexTypeMap, MapIndex: "value"},
					Index{Type: IndexTypeList, ListIndex: 0},
					Index{Type: IndexTypeMap, MapIndex: "certificate"},
					Index{Type: IndexTypeMap, MapIndex: "private_key_pem"},
				},
			}))
		})
	})

	Context("GetPropertiesIndexMap", func() {
		var (
			productProperties map[string]interface{}
			indexMap          IndexMap
			err               error
		)

		JustBeforeEach(func() {
			indexMap = GetPropertiesIndexMap(productProperties)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("When the index keys are placeholders", func() {
			Context("when the product properties contain a nested map", func() {
				BeforeEach(func() {
					productProperties = map[string]interface{}{
						".uaa.service_provider_key_credentials": map[string]interface{}{
							"value": map[string]interface{}{
								"cert_pem":        "((uaa/service_provider_key_credentials/certificate))",
								"private_key_pem": "((uaa/service_provider_key_credentials/privatekey))",
							},
						},
					}
				})

				It("creates a valid IndexMap", func() {
					Expect(indexMap).To(Equal(IndexMap{
						"((uaa/service_provider_key_credentials/certificate))": []Index{
							Index{Type: IndexTypeMap, MapIndex: ".uaa.service_provider_key_credentials"},
							Index{Type: IndexTypeMap, MapIndex: "value"},
							Index{Type: IndexTypeMap, MapIndex: "cert_pem"},
						},
						"((uaa/service_provider_key_credentials/privatekey))": []Index{
							Index{Type: IndexTypeMap, MapIndex: ".uaa.service_provider_key_credentials"},
							Index{Type: IndexTypeMap, MapIndex: "value"},
							Index{Type: IndexTypeMap, MapIndex: "private_key_pem"},
						},
					}))
				})
			})

			Context("when the product properties contain a nested slice", func() {
				BeforeEach(func() {
					productProperties = map[string]interface{}{
						".properties.networking_poe_ssl_certs": map[string]interface{}{
							"value": []interface{}{
								map[string]interface{}{
									"certificate": map[string]interface{}{
										"cert_pem":        "((networking_poe_ssl_certs_0/certificate))",
										"private_key_pem": "((networking_poe_ssl_certs_0/privatekey))",
									},
									"name": "((networking_poe_ssl_certs_0/name))",
								},
							},
						},
					}
				})

				It("creates a valid IndexMap", func() {
					Expect(indexMap).To(Equal(IndexMap{
						"((networking_poe_ssl_certs_0/certificate))": []Index{
							Index{Type: IndexTypeMap, MapIndex: ".properties.networking_poe_ssl_certs"},
							Index{Type: IndexTypeMap, MapIndex: "value"},
							Index{Type: IndexTypeList, ListIndex: 0},
							Index{Type: IndexTypeMap, MapIndex: "certificate"},
							Index{Type: IndexTypeMap, MapIndex: "cert_pem"},
						},
						"((networking_poe_ssl_certs_0/privatekey))": []Index{
							Index{Type: IndexTypeMap, MapIndex: ".properties.networking_poe_ssl_certs"},
							Index{Type: IndexTypeMap, MapIndex: "value"},
							Index{Type: IndexTypeList, ListIndex: 0},
							Index{Type: IndexTypeMap, MapIndex: "certificate"},
							Index{Type: IndexTypeMap, MapIndex: "private_key_pem"},
						},
						"((networking_poe_ssl_certs_0/name))": []Index{
							Index{Type: IndexTypeMap, MapIndex: ".properties.networking_poe_ssl_certs"},
							Index{Type: IndexTypeMap, MapIndex: "value"},
							Index{Type: IndexTypeList, ListIndex: 0},
							Index{Type: IndexTypeMap, MapIndex: "name"},
						},
					}))
				})
			})
		})

		Context("when the index keys are hardcoded values", func() {
			BeforeEach(func() {
				productProperties = map[string]interface{}{
					".properties.networking_poe_ssl_certs": map[string]interface{}{
						"value": []interface{}{
							map[string]interface{}{
								"certificate": map[string]interface{}{
									"cert_pem":        "cert key",
									"private_key_pem": "***",
								},
								"name": "bob",
							},
						},
					},
				}
			})

			It("creates a valid IndexMap", func() {
				Expect(indexMap).To(Equal(IndexMap{
					"cert key": []Index{
						Index{Type: IndexTypeMap, MapIndex: ".properties.networking_poe_ssl_certs"},
						Index{Type: IndexTypeMap, MapIndex: "value"},
						Index{Type: IndexTypeList, ListIndex: 0},
						Index{Type: IndexTypeMap, MapIndex: "certificate"},
						Index{Type: IndexTypeMap, MapIndex: "cert_pem"},
					},
					"***": []Index{
						Index{Type: IndexTypeMap, MapIndex: ".properties.networking_poe_ssl_certs"},
						Index{Type: IndexTypeMap, MapIndex: "value"},
						Index{Type: IndexTypeList, ListIndex: 0},
						Index{Type: IndexTypeMap, MapIndex: "certificate"},
						Index{Type: IndexTypeMap, MapIndex: "private_key_pem"},
					},
					"bob": []Index{
						Index{Type: IndexTypeMap, MapIndex: ".properties.networking_poe_ssl_certs"},
						Index{Type: IndexTypeMap, MapIndex: "value"},
						Index{Type: IndexTypeList, ListIndex: 0},
						Index{Type: IndexTypeMap, MapIndex: "name"},
					},
				}))
			})
		})
	})

	Context("LookupPropertyWithIndexList", func() {
		var (
			productProperties interface{}
			indexList         []Index
			value             interface{}
			err               error
		)

		JustBeforeEach(func() {
			value, err = LookupPropertyWithIndexList(indexList, productProperties)
		})

		Context("when the indexlist doesn't reach a value", func() {
			BeforeEach(func() {
				indexList = []Index{
					Index{Type: IndexTypeMap, MapIndex: ".uaa.service_provider_key_credentials"},
					Index{Type: IndexTypeMap, MapIndex: "value"},
				}
				productProperties = map[string]interface{}{
					".uaa.service_provider_key_credentials": map[string]interface{}{
						"value": map[string]interface{}{
							"cert_pem":        "-----BEGIN CERTIFICATE-----\nbeep\n-----END CERTIFICATE-----\n",
							"private_key_pem": "***",
						},
					},
				}
			})

			It("returns a NoValueAtEndOfIndex", func() {
				Expect(err).To(MatchError(NoValueAtEndOfIndexError{
					RemainingValue: map[string]interface{}{
						"cert_pem":        "-----BEGIN CERTIFICATE-----\nbeep\n-----END CERTIFICATE-----\n",
						"private_key_pem": "***",
					},
				}))
			})
		})

		Context("when there is a mismatched type in the indexlist", func() {
			Context("when the mismatched provided type is a Map", func() {
				BeforeEach(func() {
					indexList = []Index{
						Index{Type: IndexTypeMap, MapIndex: ".uaa.service_provider_key_credentials"},
						Index{Type: IndexTypeMap, MapIndex: "value"},
						Index{Type: IndexTypeMap, MapIndex: "cert_pem"},
					}

					productProperties = map[string]interface{}{
						".uaa.service_provider_key_credentials": map[string]interface{}{
							"value": []interface{}{
								map[string]interface{}{
									"cert_pem":        "-----BEGIN CERTIFICATE-----\nbeep\n-----END CERTIFICATE-----\n",
									"private_key_pem": "***",
								},
							},
						},
					}
				})

				It("returns an InvalidIndexTypeError", func() {
					Expect(err).To(MatchError(InvalidIndexTypeError{
						ProvidedType: IndexTypeMap,
						RequiredType: IndexTypeList,
					}))
				})
			})

			Context("when the mismatched provided type is a List", func() {
				BeforeEach(func() {
					indexList = []Index{
						Index{Type: IndexTypeMap, MapIndex: ".uaa.service_provider_key_credentials"},
						Index{Type: IndexTypeList, ListIndex: 1},
						Index{Type: IndexTypeMap, MapIndex: "cert_pem"},
					}

					productProperties = map[string]interface{}{
						".uaa.service_provider_key_credentials": map[string]interface{}{
							"value": map[string]interface{}{
								"cert_pem":        "-----BEGIN CERTIFICATE-----\nbeep\n-----END CERTIFICATE-----\n",
								"private_key_pem": "***",
							},
						},
					}
				})

				It("returns an InvalidIndexTypeError", func() {
					Expect(err).To(MatchError(InvalidIndexTypeError{
						ProvidedType: IndexTypeList,
						RequiredType: IndexTypeMap,
					}))
				})
			})
		})

		Context("when the product properties contain a nested map", func() {
			BeforeEach(func() {
				indexList = []Index{
					Index{Type: IndexTypeMap, MapIndex: ".uaa.service_provider_key_credentials"},
					Index{Type: IndexTypeMap, MapIndex: "value"},
					Index{Type: IndexTypeMap, MapIndex: "cert_pem"},
				}
				productProperties = map[string]interface{}{
					".uaa.service_provider_key_credentials": map[string]interface{}{
						"value": map[string]interface{}{
							"cert_pem":        "-----BEGIN CERTIFICATE-----\nbeep\n-----END CERTIFICATE-----\n",
							"private_key_pem": "***",
						},
					},
				}
			})

			It("creates a valid IndexMap", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(value).To(Equal("-----BEGIN CERTIFICATE-----\nbeep\n-----END CERTIFICATE-----\n"))
			})
		})

		// Context("when the index map skips an extra 'value' nesting", func() {
		// 	BeforeEach(func() {
		// 		indexList = []Index{
		// 			Index{Type: IndexTypeMap, MapIndex: ".uaa.service_provider_key_credentials"},
		// 			Index{Type: IndexTypeMap, MapIndex: "cert_pem"},
		// 		}
		// 		productProperties = map[string]interface{}{
		// 			".uaa.service_provider_key_credentials": map[string]interface{}{
		// 				"value": map[string]interface{}{
		// 					"cert_pem":        "-----BEGIN CERTIFICATE-----\nbeep\n-----END CERTIFICATE-----\n",
		// 					"private_key_pem": "***",
		// 				},
		// 			},
		// 		}
		// 	})

		// 	It("creates a valid IndexMap", func() {
		// 		Expect(value).To(Equal("-----BEGIN CERTIFICATE-----\nbeep\n-----END CERTIFICATE-----\n"))
		// 	})
		// })

		Context("when the product properties contain a nested list", func() {
			BeforeEach(func() {
				indexList = []Index{
					Index{Type: IndexTypeMap, MapIndex: ".properties.networking_poe_ssl_certs"},
					Index{Type: IndexTypeMap, MapIndex: "value"},
					Index{Type: IndexTypeList, ListIndex: 0},
					Index{Type: IndexTypeMap, MapIndex: "certificate"},
					Index{Type: IndexTypeMap, MapIndex: "cert_pem"},
				}
				productProperties = map[string]interface{}{
					".properties.networking_poe_ssl_certs": map[string]interface{}{
						"value": []interface{}{
							map[string]interface{}{
								"certificate": map[string]interface{}{
									"cert_pem":        "-----BEGIN CERTIFICATE-----\nbeep\n-----END CERTIFICATE-----\n",
									"private_key_pem": "***",
								},
								"name": "((networking_poe_ssl_certs_0/name))",
							},
						},
					},
				}
			})

			It("creates a valid IndexMap", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(value).To(Equal("-----BEGIN CERTIFICATE-----\nbeep\n-----END CERTIFICATE-----\n"))
			})
		})
	})
})
