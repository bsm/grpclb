require 'bundler/setup'
require 'bundler/gem_tasks'

require 'rubocop/rake_task'
RuboCop::RakeTask.new(:rubocop)

# Release is disabled as it tries to create a tag on github.
# Run `rake build` and `gem push` instead.
Rake::Task[:release].clear

task default: %i[rubocop]
