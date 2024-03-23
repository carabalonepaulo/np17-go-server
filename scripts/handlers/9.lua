local listener = require 'listener'

return function(_, message)
  listener:send_to_all(string.format('<9>%s</9>', message))
end
