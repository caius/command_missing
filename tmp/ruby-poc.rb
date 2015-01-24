#!/usr/bin/env ruby

def paths
  @paths ||= ENV["PATH"].split(":")
end

def cmd_path(cmd)
  paths.map { |path| potential = File.join(path, cmd); potential if File.exist?(potential) }.compact.first
end

def candidates(array)
  array[0] ||= ""
  array[1] ||= ""
  array[0].size.downto(1).map { |x| array[0][0, x] } | \
  (1..array[1].size).map { |x| array[0] + array[1][0, x] }
end

if found_it = candidates(ARGV).map { |cmd| cmd_path(cmd) }.compact.first
  puts found_it
  exit(0)
else
  exit(1)
end
