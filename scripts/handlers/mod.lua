local listener = require 'listener'

return function(sender_id)
  local motd = 'message of the day'
  listener:send_to(sender_id, string.format('<mod>%s</mod>', motd))
end
