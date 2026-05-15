return {
	"nvim-treesitter/nvim-treesitter",
	branch = "main",

	build = ":TSUpdate",

	config = function()
		local parsers = {
			"bash",
			"c",
			"diff",
			"cpp",
			"html",
			"lua",
			"luadoc",
			"markdown",
			"markdown_inline",
			"query",
			"vim",
			"vimdoc",
			"typescript",
			"javascript",
			"go",
			"php",
		}

		require("nvim-treesitter").install(parsers)

		vim.api.nvim_create_autocmd("FileType", {
			callback = function(args)
				local buf = args.buf
				local ft = args.match

				local lang = vim.treesitter.language.get_lang(ft)
				if not lang then
					return
				end

				pcall(vim.treesitter.start, buf, lang)
			end,
		})
	end,
}
