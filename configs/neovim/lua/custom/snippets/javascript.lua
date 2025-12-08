local ls = require("luasnip")
local rep = require("luasnip.extras").rep
local ts_utils = require("nvim-treesitter.ts_utils")

local s = ls.snippet
local t = ls.text_node
local i = ls.insert_node
local f = ls.function_node

local filename = function()
	return vim.fn.expand("%:t")
end

local current_function_name = function()
	local node = ts_utils.get_node_at_cursor()
	while node do
		if node:type():match("function") then
			-- Get the text of the function name node
			local name_node = node:field("name")[1]
			if name_node then
				return vim.treesitter.get_node_text(name_node, 0)
			end
		end
		node = node:parent()
	end
	return "function"
end

local js_snippets = {
	s("dlog", {
		t("console.log('["),
		f(filename),
		t("] "),
		f(current_function_name),
		t(":"),
		i(1, "variable"),
		t(" -> ', "),
		rep(1),
		t(");"),
	}),
}

ls.add_snippets("javascript", js_snippets)
ls.add_snippets("typescript", js_snippets)
