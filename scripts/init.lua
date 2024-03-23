local printf = require 'printf'
local listener = require 'listener'
local log = require 'log'
log.printf = function(...) log.print(string.format(...)) end

return {
  init = function()
  end,

  deinit = function()
  end,

  update = function()
  end,

  on_client_connected = function(client_id)
    log.printf('Client `%d` connected!', client_id)
  end,

  on_client_disconnected = function(client_id)
    log.printf('Client `%d` disconnected!', client_id)
  end,

  on_data_received = function(client_id, message_name, message_content)
    local success, handler = pcall(require, 'handlers.' .. message_name)
    if success then
      handler(client_id, message_content)
    else
      log.printf('No handler found for packet `%s`!', message_name)
    end
  end,
}
