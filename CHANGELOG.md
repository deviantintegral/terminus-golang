# Changelog

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
