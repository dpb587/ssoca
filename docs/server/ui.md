# Frontend UI

A simple UI is included in the [`server/ui`](../../server/ui) directory to support users directly browsing the the ssoca endpoint. It can be configured with the [`docroot`](../service/docroot) service and supports showing initial configuration steps, binary downloads, and minor customizations.


## Parameters

The following optional settings may be configured in the `env.metadata` to affect the appearance.

 * `ui.color` - a [known](http://tachyons.io/docs/themes/skins/) color name for the background of pages
 * `ui.link.(title)` - URL(s) to external documentation or resources which may be useful to the user
 * `ui.usage` - simple snippets documenting common commands which may be useful to the user

Web pages also reference the `name`, `title`, and `banner` properties of the `env` settings.
