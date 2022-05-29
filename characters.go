package aegirdungeons

const (
	Unknown = iota
	Templar
	Berserker
	Shadowblade
	Purifier
	Guardian
	Archon
	Marauder
	Siegebreaker
	Inquisitor
	Seer
	Sage
	Raider
	Mini
	Corsair
	Outlaw
	Blademaster
)

var TestnetAsaToClass = map[uint64]uint64{
	21849216: Templar,
	46001022: Berserker,
	73723335: Mini,
}

var TestnetAsaToEmoji = map[uint64]string{
	21849216: "<:gregor:911337902078840853>",
	46001022: "<:hjalmar:911339768187605002>",
	73723335: "<:miniaowl:940384267190546442>",
}

var MainnetAsaToClass = map[uint64]uint64{
	382606747: Templar,
	404079915: Berserker,
	410746800: Shadowblade,
	503250857: Shadowblade,
	423590752: Purifier,
	557904219: Purifier,
	436620559: Guardian,
	447390956: Archon,
	459238723: Marauder,
	469434825: Siegebreaker,
	511246825: Inquisitor,
	525573439: Seer,
	544736887: Sage,
	562258870: Raider,
	573509824: Corsair,
	587208857: Outlaw,
	596914358: Guardian,
	601824851: Blademaster,

	// Algomonz
	757517034: Archon,
	757521685: Raider,

	// AOWL MINIS
	577370701: Mini,
	577381222: Mini,
	577381383: Mini,
	577382314: Mini,
	577382998: Mini,
	577385856: Mini,
	577459363: Mini,
	577462356: Mini,
	577463437: Mini,
	577464449: Mini,
	577464950: Mini,
	577465711: Mini,
	577466295: Mini,
	577466830: Mini,
	577467279: Mini,

	// Arcadian Tales MINIS
	607980607: Mini,
	607984641: Mini,
	607984735: Mini,
	607984886: Mini,
	607985003: Mini,
	607985168: Mini,
	607985294: Mini,
	607985431: Mini,
	607985566: Mini,
	607985660: Mini,

	// Raven MINIS
	653332430: Mini,
	653332824: Mini,
	653332987: Mini,
	653334069: Mini,
	653334675: Mini,
	653334906: Mini,
	653335044: Mini,
	653335776: Mini,
	653335839: Mini,
	653335901: Mini,
	653336647: Mini,
	653337226: Mini,
	653337479: Mini,
	653337621: Mini,
	653337761: Mini,
	653338146: Mini,
	653338536: Mini,
	653338843: Mini,
	653339311: Mini,
	653340148: Mini,
	653340742: Mini,
	653341000: Mini,
	653341152: Mini,
	653341285: Mini,
	653341476: Mini,
	653341630: Mini,
	653341762: Mini,
	653341920: Mini,
	653342093: Mini,
	653342361: Mini,
	653342667: Mini,
	653342977: Mini,
	653343119: Mini,
	653343260: Mini,
	653343648: Mini,
	653343983: Mini,
	653344114: Mini,
	653344703: Mini,
	653345018: Mini,
	653345136: Mini,
	653345593: Mini,
	653345741: Mini,
	653345900: Mini,
	653346409: Mini,
	653346609: Mini,
	653346712: Mini,
	653347184: Mini,
	653347288: Mini,
	653347731: Mini,
	653348381: Mini,
	653348515: Mini,
	653349837: Mini,
	653350169: Mini,
	653350246: Mini,
	653350623: Mini,
	653350732: Mini,
	653350841: Mini,
	653350956: Mini,
	653351101: Mini,
	653351332: Mini,
	653351420: Mini,
	653351586: Mini,
	653351699: Mini,
	653353351: Mini,
	653353480: Mini,
	653353654: Mini,
	653353837: Mini,
	653354307: Mini,
	653354458: Mini,
	653354596: Mini,
	653354926: Mini,
	653355203: Mini,
	653355483: Mini,
	653355752: Mini,
	653355863: Mini,
	653356140: Mini,
	653356297: Mini,
	653356596: Mini,
	653356740: Mini,

	// Nack Mache
	708382784: Mini,
	708403399: Mini,
	708410183: Mini,
	708413631: Mini,
	708415128: Mini,
	708416381: Mini,
	708429855: Mini,
	708430866: Mini,
	708431583: Mini,
	708432062: Mini,
	708432757: Mini,
	708433235: Mini,
	708433936: Mini,
	708434736: Mini,
	708437171: Mini,
	708439668: Mini,
	708440857: Mini,
	708441949: Mini,
	708443428: Mini,
	708444258: Mini,
}

var MainnetAsaToEmoji = map[uint64]string{
	382606747: "<:gregor:911337902078840853>",
	404079915: "<:hjalmar:911339768187605002>",
	410746800: "<:sigurd:911339223263629342>",
	423590752: "<:valeria:911277491027587143>",
	436620559: "<:kor:912405242623193128>",
	447390956: "<:nadia:914936931127791636>",
	459238723: "<:jargal:917555878813650974>",
	469434825: "<:herja:921075549835784283>",
	503250857: "<:shadawolblade:925921171998920776>",
	511246825: "<:ivan:925917223405637672>",
	525573439: "<:orm:928423348298453072>",
	544736887: "<:yara:933107248727748620>",
	562258870: "<:nazir:934891652089270314>",
	557904219: "<:arcadianvaleria:935205763654246402>",
	573509824: "<:rashad:938098042412879893>",
	757517034: "<:algomonz_nadia:980464036296687666>",
	757521685: "<:algomonz_nazir:980463480794673193>",

	577370701: "<:miniaowl:940384267190546442>",
	577381222: "<:miniaowl:940384267190546442>",
	577381383: "<:miniaowl:940384267190546442>",
	577382314: "<:miniaowl:940384267190546442>",
	577382998: "<:miniaowl:940384267190546442>",
	577385856: "<:miniaowl:940384267190546442>",
	577459363: "<:miniaowl:940384267190546442>",
	577462356: "<:miniaowl:940384267190546442>",
	577463437: "<:miniaowl:940384267190546442>",
	577464449: "<:miniaowl:940384267190546442>",
	577464950: "<:miniaowl:940384267190546442>",
	577465711: "<:miniaowl:940384267190546442>",
	577466295: "<:miniaowl:940384267190546442>",
	577466830: "<:miniaowl:940384267190546442>",
	577467279: "<:miniaowl:940384267190546442>",
	587208857: "<:kurt:941032362131750922>",
	596914358: "<:akita:941186702704250960>",
	601824851: "<:uma:943925053316284466>",
	607980607: "<:mini_arcadian_dog:951673127031689316>",
	607984641: "<:mini_arcadian_dog:951673127031689316>",
	607984735: "<:mini_arcadian_dog:951673127031689316>",
	607984886: "<:mini_arcadian_dog:951673127031689316>",
	607985003: "<:mini_arcadian_dog:951673127031689316>",
	607985168: "<:mini_arcadian_dragon:951673282690678868>",
	607985294: "<:mini_arcadian_dragon:951673282690678868>",
	607985431: "<:mini_arcadian_dragon:951673282690678868>",
	607985566: "<:mini_arcadian_dragon:951673282690678868>",
	607985660: "<:mini_arcadian_dragon:951673282690678868>",
	653332430: "<:mini_raven:952005359978020864>",
	653332824: "<:mini_raven:952005359978020864>",
	653332987: "<:mini_raven:952005359978020864>",
	653334069: "<:mini_raven:952005359978020864>",
	653334675: "<:mini_raven:952005359978020864>",
	653334906: "<:mini_raven:952005359978020864>",
	653335044: "<:mini_raven:952005359978020864>",
	653335776: "<:mini_raven:952005359978020864>",
	653335839: "<:mini_raven:952005359978020864>",
	653335901: "<:mini_raven:952005359978020864>",
	653336647: "<:mini_raven:952005359978020864>",
	653337226: "<:mini_raven:952005359978020864>",
	653337479: "<:mini_raven:952005359978020864>",
	653337621: "<:mini_raven:952005359978020864>",
	653337761: "<:mini_raven:952005359978020864>",
	653338146: "<:mini_raven:952005359978020864>",
	653338536: "<:mini_raven:952005359978020864>",
	653338843: "<:mini_raven:952005359978020864>",
	653339311: "<:mini_raven:952005359978020864>",
	653340148: "<:mini_raven:952005359978020864>",
	653340742: "<:mini_raven:952005359978020864>",
	653341000: "<:mini_raven:952005359978020864>",
	653341152: "<:mini_raven:952005359978020864>",
	653341285: "<:mini_raven:952005359978020864>",
	653341476: "<:mini_raven:952005359978020864>",
	653341630: "<:mini_raven:952005359978020864>",
	653341762: "<:mini_raven:952005359978020864>",
	653341920: "<:mini_raven:952005359978020864>",
	653342093: "<:mini_raven:952005359978020864>",
	653342361: "<:mini_raven:952005359978020864>",
	653342667: "<:mini_raven:952005359978020864>",
	653342977: "<:mini_raven:952005359978020864>",
	653343119: "<:mini_raven:952005359978020864>",
	653343260: "<:mini_raven:952005359978020864>",
	653343648: "<:mini_raven:952005359978020864>",
	653343983: "<:mini_raven:952005359978020864>",
	653344114: "<:mini_raven:952005359978020864>",
	653344703: "<:mini_raven:952005359978020864>",
	653345018: "<:mini_raven:952005359978020864>",
	653345136: "<:mini_raven:952005359978020864>",
	653345593: "<:mini_raven:952005359978020864>",
	653345741: "<:mini_raven:952005359978020864>",
	653345900: "<:mini_raven:952005359978020864>",
	653346409: "<:mini_raven:952005359978020864>",
	653346609: "<:mini_raven:952005359978020864>",
	653346712: "<:mini_raven:952005359978020864>",
	653347184: "<:mini_raven:952005359978020864>",
	653347288: "<:mini_raven:952005359978020864>",
	653347731: "<:mini_raven:952005359978020864>",
	653348381: "<:mini_raven:952005359978020864>",
	653348515: "<:mini_raven:952005359978020864>",
	653349837: "<:mini_raven:952005359978020864>",
	653350169: "<:mini_raven:952005359978020864>",
	653350246: "<:mini_raven:952005359978020864>",
	653350623: "<:mini_raven:952005359978020864>",
	653350732: "<:mini_raven:952005359978020864>",
	653350841: "<:mini_raven:952005359978020864>",
	653350956: "<:mini_raven:952005359978020864>",
	653351101: "<:mini_raven:952005359978020864>",
	653351332: "<:mini_raven:952005359978020864>",
	653351420: "<:mini_raven:952005359978020864>",
	653351586: "<:mini_raven:952005359978020864>",
	653351699: "<:mini_raven:952005359978020864>",
	653353351: "<:mini_raven:952005359978020864>",
	653353480: "<:mini_raven:952005359978020864>",
	653353654: "<:mini_raven:952005359978020864>",
	653353837: "<:mini_raven:952005359978020864>",
	653354307: "<:mini_raven:952005359978020864>",
	653354458: "<:mini_raven:952005359978020864>",
	653354596: "<:mini_raven:952005359978020864>",
	653354926: "<:mini_raven:952005359978020864>",
	653355203: "<:mini_raven:952005359978020864>",
	653355483: "<:mini_raven:952005359978020864>",
	653355752: "<:mini_raven:952005359978020864>",
	653355863: "<:mini_raven:952005359978020864>",
	653356140: "<:mini_raven:952005359978020864>",
	653356297: "<:mini_raven:952005359978020864>",
	653356596: "<:mini_raven:952005359978020864>",
	653356740: "<:mini_raven:952005359978020864>",

	708382784: "<:mini_nack_mache:974106946141556776>",
	708403399: "<:mini_nack_mache:974106946141556776>",
	708410183: "<:mini_nack_mache:974106946141556776>",
	708413631: "<:mini_nack_mache:974106946141556776>",
	708415128: "<:mini_nack_mache:974106946141556776>",
	708416381: "<:mini_nack_mache:974106946141556776>",
	708429855: "<:mini_nack_mache:974106946141556776>",
	708430866: "<:mini_nack_mache:974106946141556776>",
	708431583: "<:mini_nack_mache:974106946141556776>",
	708432062: "<:mini_nack_mache:974106946141556776>",
	708432757: "<:mini_nack_mache:974106946141556776>",
	708433235: "<:mini_nack_mache:974106946141556776>",
	708433936: "<:mini_nack_mache:974106946141556776>",
	708434736: "<:mini_nack_mache:974106946141556776>",
	708437171: "<:mini_nack_mache:974106946141556776>",
	708439668: "<:mini_nack_mache:974106946141556776>",
	708440857: "<:mini_nack_mache:974106946141556776>",
	708441949: "<:mini_nack_mache:974106946141556776>",
	708443428: "<:mini_nack_mache:974106946141556776>",
	708444258: "<:mini_nack_mache:974106946141556776>",
}
