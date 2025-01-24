-- Custom ***
-- NOTE: Remove any other theme
return {
  'ellisonleao/gruvbox.nvim',
  priority = 1000,
  config = function()
    -- Set the background and colorscheme after the plugin is loaded
    vim.o.background = 'dark'
    vim.cmd [[colorscheme gruvbox]]
  end,
}
-- ***
