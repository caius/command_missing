ARGV = %w(bre wup) # !> already initialized constant ARGV
ENV["PATH"] = "/Users/caius/go/bin:/Users/caius/bin:/Users/caius/.cabal/bin:/usr/local/sbin:/Users/caius/.gem/ruby/2.1.5/bin:/Users/caius/.rubies/ruby-2.1.5/lib/ruby/gems/2.1.0/bin:/Users/caius/.rubies/ruby-2.1.5/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:/opt/X11/bin"

def paths
  @paths ||= ENV["PATH"].split(":")
end

def cmd?(cmd)
  paths.any? { |path| File.exist? File.join(path, cmd) }
end

def candidates(array)
  (1..array[0].size).map { |x| array[0][0, x] } | \
  (1..array[1].size).map { |x| array[0] + array[1][0, x] }
end

candidates(ARGV).find { |cmd| cmd?(cmd) } || exit(0)
