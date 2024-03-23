local listener = require 'listener'

return function(sender_id)
  local version = '1.7'
  listener:send_to(sender_id, string.format('<ver>%s</ver>', version))
end
