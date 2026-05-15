local M = {}

---@param repo string
---@return string
function M.gh(repo)
	return "https://github.com/" .. repo
end

return M
