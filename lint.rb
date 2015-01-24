#!/usr/bin/env ruby

unless ARGV.size == 1 && ($binary = File.expand_path(ARGV.first)) && File.exist?($binary)
  puts "ERROR: couldn't find binary ./#{$binary.inspect}"
  exit 1
end

def assert_equal(expected, actual)
  if expected == actual
    print "."
  else
    puts "FAIL: expected #{expected.inspect}, got #{actual.inspect}"
  end
end
at_exit { puts }

output = %x[#{$binary} bash something]
assert_equal 0, $?.exitstatus
assert_equal "/bin/bash\n", output

output = %x[#{$binary} bas h something]
assert_equal 0, $?.exitstatus
assert_equal "/bin/bash\n", output

output = %x[#{$binary} bas fuck something]
assert_equal 1, $?.exitstatus
assert_equal "", output

output = %x[#{$binary} invalid something]
assert_equal 1, $?.exitstatus
assert_equal "", output
