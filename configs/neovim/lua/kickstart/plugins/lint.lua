local gh = require("kickstart.utils").gh

vim.pack.add({ gh("mfussenegger/nvim-lint") })

local lint = require("lint")
lint.linters_by_ft = {
	markdown = { "markdownlint" },
	-- [devgita]
	python = { "flake8" },
	go = { "golangcilint" },
	-- [/devgita]
}

local lint_augroup = vim.api.nvim_create_augroup("lint", { clear = true })
vim.api.nvim_create_autocmd({ "BufWritePost" }, {
	group = lint_augroup,
	callback = function()
		if vim.opt_local.modifiable:get() then
			lint.try_lint()
		end
	end,
})
