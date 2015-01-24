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

def run(args, &check_block)
  print "args: #{args.inspect} -- "
  output = %x[#{$binary} #{args.join(" ")}]
  yield($?.exitstatus, output)
  puts
end

# We don't expect to match valid binary names without any futzing - ZSH should've found them already
run(%w(bash)) do |status, output|
  assert_equal 1, status
  assert_equal "", output
end

run(%w(bas h something)) do |status, output|
  assert_equal 0, status
  assert_equal "/bin/bash something\n", output
end

run(%w(bas hsomething)) do |status, output|
  assert_equal 0, status
  assert_equal "/bin/bash something\n", output
end

run(%w(bashs omething)) do |status, output|
  assert_equal 0, status
  assert_equal "/bin/bash something\n", output
end

run(%w(bas fuck something)) do |status, output|
  assert_equal 1, status
  assert_equal "", output
end

run(%w(invalid something)) do |status, output|
  assert_equal 1, status
  assert_equal "", output
end
