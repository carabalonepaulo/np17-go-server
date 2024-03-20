local printf = function(...) print(string.format(...)) end
local listener = require 'listener'

return {
  init = function()
    -- print('LuaApi initialized!')
    -- test:send_to(0, '<0>hello</0>')
  end,

  deinit = function()
    -- print('LuaApi finalized!')
  end,

  on_client_connected = function(client_id)
    printf('Client `%d` connected!', client_id)
  end,

  on_client_disconnected = function(client_id)
    printf('Client `%d` disconnected!', client_id)
  end,

  on_data_received = function(client_id, message_name, message_content)
    printf('Message `%s` received from client `%d`!', message_name, client_id)
    -- printf('Total received: `%d`', listener:get_total_received(client_id))

    if message_name == '<0>' then
      local packet = string.format("<0 %d>'e' n=%s</0>", client_id, 'server_name')
      listener:send_to(client_id, packet)
    end

    -- listener:send_to_many('<0>hello world</0>', function(id)
    --   return false
    -- end)
  end,
}
