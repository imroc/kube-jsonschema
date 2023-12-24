# kube-jsonschema

## Usage

Neovim configuration:

```lua
--  skip lunarvim's built-in yamlls configuration if you use lunarvim:
--  vim.list_extend(lvim.lsp.automatic_configuration.skipped_servers, { "yamlls" })

local schemas = {
  -- ignore yamlls's built-in k8s json schema
  kubernetes = "",
  ["https://raw.githubusercontent.com/imroc/kube-jsonschema/master/schemas/all.json"] = {"*.yaml", "*.yml"},
}
-- optionnal but useful if you want other commonly used schemas form schemastore.
-- (The prerequisite is that the schema plugin is installed: https://github.com/b0o/schemastore.nvim)
schemas = vim.tbl_extend("force", schemas, require('schemastore').yaml.schemas())
require('lspconfig').yamlls.setup {
  settings = {
    yaml = {
      schemaStore = {
        -- You must disable built-in schemaStore support if you want to use
        -- this plugin and its advanced options like `ignore`.
        enable = false,
        -- Avoid TypeError: Cannot read properties of undefined (reading 'length')
        url = "",
      },
      schemas = schemas,
    },
  },
}
```
