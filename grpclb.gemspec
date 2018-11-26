Gem::Specification.new do |s|
  s.name          = 'grpclb'
  s.summary       = 'grpclb ruby protocol'
  s.version       = File.read(File.expand_path('.version', __dir__)).strip
  s.authors       = ['Black Square Media']
  s.platform      = Gem::Platform::RUBY
  s.files         = `git ls-files ruby`.split("\n")
  s.require_paths = ['ruby/lib']

  s.add_runtime_dependency 'grpc'
  s.add_development_dependency 'grpc-tools'
  s.add_development_dependency 'rake'
  s.add_development_dependency 'rubocop'
end
