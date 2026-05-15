local gh = require("kickstart.utils").gh

vim.pack.add({ gh("NMAC427/guess-indent.nvim") })
require("guess-indent").setup({})
