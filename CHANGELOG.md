# Changelog

## [0.5.0](https://github.com/deviantintegral/terminus-golang/compare/v0.4.0...v0.5.0) (2025-11-28)


### Features

* add owner, team, upstream filters to site:list ([da0d54f](https://github.com/deviantintegral/terminus-golang/commit/da0d54fb058600aa3c77c24887f3ac7dbcad56e8))
* add serialization for org:list display ([c72ba28](https://github.com/deviantintegral/terminus-golang/commit/c72ba286ad2d22a190633e98950cebab80c1d1e1))
* add Serializer interface for PHP-like CSV ([d609317](https://github.com/deviantintegral/terminus-golang/commit/d609317dbf203464fb2ce29999660ace54a0151d))
* mark site:list command as tested ([#101](https://github.com/deviantintegral/terminus-golang/issues/101)) ([c37bb44](https://github.com/deviantintegral/terminus-golang/commit/c37bb441a1f7310eb4f75f266d20085a7cd096fe))
* update CSV column names to match PHP format ([3a60ab4](https://github.com/deviantintegral/terminus-golang/commit/3a60ab4bd59cc893dc42783ed45b217267ca5a52))
* update Organization and Profile models to match API ([7fc67c6](https://github.com/deviantintegral/terminus-golang/commit/7fc67c61ad0141e5d9863c50756f4391cf78fdb5))
* use friendly plan names from API ([c65056b](https://github.com/deviantintegral/terminus-golang/commit/c65056b8451af1ae3249c48ef6a8577c7cb6da0b))
* use friendly region names from API ([649f8e6](https://github.com/deviantintegral/terminus-golang/commit/649f8e6a57addff3c1b54ac0615efc0ffb8a3d4f))


### Bug Fixes

* add GOPATH/bin to PATH in SessionStart hook ([58725be](https://github.com/deviantintegral/terminus-golang/commit/58725be4c9f0bcfa44c0bd4d438c71c40f386160))
* improve org name resolution error messages ([8d8340b](https://github.com/deviantintegral/terminus-golang/commit/8d8340bb1e6f7925962cb39000b3188de68829ee))
* populate Memberships field in site:list ([720f17e](https://github.com/deviantintegral/terminus-golang/commit/720f17e378c9b3f1034af8e0c45f4c28777fa5e3))
* resolve org name to UUID in site:list ([80e745f](https://github.com/deviantintegral/terminus-golang/commit/80e745f79a9614d3a1159a95159b9e56cc25a7d2))
* resolve Serializer interface not being used for CSV ([eccf35e](https://github.com/deviantintegral/terminus-golang/commit/eccf35eb8b48175743236310b41862ac2a592173))
* serialize logo width and height strings into ints ([3e249ec](https://github.com/deviantintegral/terminus-golang/commit/3e249ec5af6c462ef3bc77deb0fb5f6e4043c097))
* site:list --team to include memberships by team ([c8cbcfd](https://github.com/deviantintegral/terminus-golang/commit/c8cbcfd33440f18deab8b29d6d8e280f45573ffd))
* site:list --upstream takes uuids, not names ([040beb3](https://github.com/deviantintegral/terminus-golang/commit/040beb350b5b196b2d1b4b40352e6c450c543666))
* symlink goimports to /usr/local/bin ([6f7c456](https://github.com/deviantintegral/terminus-golang/commit/6f7c45603837d45f00bd0fbda1f6f3fd7752994c))

## [0.4.0](https://github.com/deviantintegral/terminus-golang/compare/v0.3.0...v0.4.0) (2025-11-21)


### Features

* add confirmation for destructive env commands ([#80](https://github.com/deviantintegral/terminus-golang/issues/80)) ([efb6235](https://github.com/deviantintegral/terminus-golang/commit/efb623571f595a9e60fa7e386531c3a465f9e786))
* add dashboard:view command ([#82](https://github.com/deviantintegral/terminus-golang/issues/82)) ([ba818da](https://github.com/deviantintegral/terminus-golang/commit/ba818da77cae9675d6de3833d2f77e27a6cc274f))
* add org memberships to site:list ([3391645](https://github.com/deviantintegral/terminus-golang/commit/339164578540c7dab36741858e5dff9b3128d5a0))
* add upstream label to site:info output ([57c1821](https://github.com/deviantintegral/terminus-golang/commit/57c1821bf9fb11ba247385fe4e19180c18828aef))
* **art:** commands have been tested and are complete ([#84](https://github.com/deviantintegral/terminus-golang/issues/84)) ([d54e728](https://github.com/deviantintegral/terminus-golang/commit/d54e728ab564e6333c1af7b1d0ee668034519d42))
* format upstream as uuid: url ([04d7786](https://github.com/deviantintegral/terminus-golang/commit/04d77861a8de18b48d51a95e6b38e629cebe723e))
* remove TRACE mode output truncation ([#81](https://github.com/deviantintegral/terminus-golang/issues/81)) ([7e2e76f](https://github.com/deviantintegral/terminus-golang/commit/7e2e76fe66a86ebef260d2322e92710a83f144b4))
* remove upstream field from site:list output ([#83](https://github.com/deviantintegral/terminus-golang/issues/83)) ([ab08f6b](https://github.com/deviantintegral/terminus-golang/commit/ab08f6b311f2dc173b93a8b3c97aa0e507abe64a))
* use args for site:create command ([1327035](https://github.com/deviantintegral/terminus-golang/commit/13270351dcb8729712f94c5bc6ec3094ea681a3d))


### Bug Fixes

* add missing newlines to JSON test fixtures ([#76](https://github.com/deviantintegral/terminus-golang/issues/76)) ([250aa83](https://github.com/deviantintegral/terminus-golang/commit/250aa83a402e9a568747660196cbebdad856e332))
* add site_state=true to site info request ([180fc98](https://github.com/deviantintegral/terminus-golang/commit/180fc9895525f0496e5c1a343a3b4bfe1f3276e1))
* **api:** add universal site name resolution helpers ([caeb7e6](https://github.com/deviantintegral/terminus-golang/commit/caeb7e68fcf8d8968c7d43fafff4b2f925e1578d))
* **api:** correct workflow types and parameters ([dd57915](https://github.com/deviantintegral/terminus-golang/commit/dd579150abd37299bca928c3cf2b4e72d8712d41))
* **api:** upstream name and site id resolution ([#77](https://github.com/deviantintegral/terminus-golang/issues/77)) ([57b9ede](https://github.com/deviantintegral/terminus-golang/commit/57b9ededcd756bae6ffc2b79dc2f6d6b69a8ef54))
* **site:** add name resolution to all site methods ([3538587](https://github.com/deviantintegral/terminus-golang/commit/35385871d9aff76c8848cb597f067a5885939a56))
* **site:** improve site:create workflow handling ([afeeede](https://github.com/deviantintegral/terminus-golang/commit/afeeede410774795c4ff980bb851d9aadbcb21d8))
* **site:** resolve site name in site:delete ([cbfff83](https://github.com/deviantintegral/terminus-golang/commit/cbfff830071e6c74bdc88df5a4345261b1a3116d))
* **site:** support name lookup for site:info ([4818031](https://github.com/deviantintegral/terminus-golang/commit/48180314643f8657c3bef8f837b5e443c50204bb))
* use cursor-based pagination in GetPaged ([ab640c3](https://github.com/deviantintegral/terminus-golang/commit/ab640c3d90d84637a7329f35feff8664ebbfc285))
* use delete_site workflow for site deletion ([f303903](https://github.com/deviantintegral/terminus-golang/commit/f303903eaf14c291ea82e0b9b77df2e75d100574))
* use user workflow endpoint for deletion ([#71](https://github.com/deviantintegral/terminus-golang/issues/71)) ([60f283c](https://github.com/deviantintegral/terminus-golang/commit/60f283ca8494a62f3eb2095a9b8bfa3e19c296ab))
* use workflows endpoints for site creation ([544335a](https://github.com/deviantintegral/terminus-golang/commit/544335a09b8409b4065ae67cd22ecc0ded3d111f))

## [0.3.0](https://github.com/deviantintegral/terminus-golang/compare/v0.2.0...v0.3.0) (2025-11-20)


### Features

* **art:** add metal art command ([17c1351](https://github.com/deviantintegral/terminus-golang/commit/17c13510bf9b3ca701d0c02617340c67f099684d))
* **auth:** add flexible login token options ([dad7abf](https://github.com/deviantintegral/terminus-golang/commit/dad7abfa13ab6ae3b4c5fc4d8a07a9408fbf9d50))
* implement all list commands ([43c0a55](https://github.com/deviantintegral/terminus-golang/commit/43c0a559bbd7ac190b681d92f4249383e755ab14))
* implement info commands ([1513020](https://github.com/deviantintegral/terminus-golang/commit/1513020bf2fe4fbafb95e17da64b97323f90687d))
* **redis:** add redis:enable and redis:disable ([406164b](https://github.com/deviantintegral/terminus-golang/commit/406164baa2601eed0b79bb275e0ad3015500ffb1))


### Bug Fixes

* **api:** add client field to login request ([daf72d8](https://github.com/deviantintegral/terminus-golang/commit/daf72d86572a3ff9fe878d083f08f9b6dcb400ed))
* **api:** return errors for HTTP 4XX/5XX responses ([d927c61](https://github.com/deviantintegral/terminus-golang/commit/d927c6130ec10b3aacd2ad23e54beaf0d733a83a))
* **art:** make the horns recognizable ([14675f2](https://github.com/deviantintegral/terminus-golang/commit/14675f26c6acbf063696fb18a461c59825989206))
* **auth:** extract raw token from PHP format ([551d28b](https://github.com/deviantintegral/terminus-golang/commit/551d28bd1a2128f370d8cba552acc37a2edb40d3))
* **auth:** flatten whoami output format ([15fd6f8](https://github.com/deviantintegral/terminus-golang/commit/15fd6f82122965ce82d22bce5d146f9f8dd957e6))
* log errors in main and use colon separators ([9c60bbb](https://github.com/deviantintegral/terminus-golang/commit/9c60bbb282c79e95976124171c3945953e133230))
* **output:** match PHP terminus table format ([dfa5083](https://github.com/deviantintegral/terminus-golang/commit/dfa50831891f91707fed5eb8f35449862b60d16e))

## [0.2.0](https://github.com/deviantintegral/terminus-golang/compare/v0.1.0...v0.2.0) (2025-11-19)


### Features

* add art commands with tests ([f7deee8](https://github.com/deviantintegral/terminus-golang/commit/f7deee8d01f60de6aafb2debea3c86430a2bacd7))


### Bug Fixes

* update auth:whoami to use /users/{id} endpoint ([0cb1ee3](https://github.com/deviantintegral/terminus-golang/commit/0cb1ee3ee0f91a7b074125b7de2b26eeb56e8ba4))
* update org and site list endpoints to match PHP Terminus implementation ([4a384d9](https://github.com/deviantintegral/terminus-golang/commit/4a384d9b5b7d82bd7f4347d50c006f54f48c7cd3))
