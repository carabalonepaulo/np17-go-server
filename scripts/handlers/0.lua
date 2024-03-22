local listener = require 'listener'

return function(sender_id, message)
  local packet = string.format("<0 %d>'e' n=%s</0>", sender_id, 'Server')
  listener:send_to(sender_id, packet)
end
