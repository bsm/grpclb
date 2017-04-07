Gem::Specification.new do |s|
	s.name          = 'grpclb'
	s.summary       = 'grpclb protocol'
	s.version       = '0.1.0'
	s.authors       = ['Black Square Media']
	s.platform      = Gem::Platform::RUBY
	s.files         = `git ls-files ruby`.split("\n")
	s.require_paths = ['ruby/lib']

	s.add_runtime_dependency 'google-protobuf'
	s.add_runtime_dependency 'grpc'
	s.add_development_dependency 'grpc-tools'
end
