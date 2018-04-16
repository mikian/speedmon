TARGETS = {
  linux: %w{mipsle mips64}
}

task :all do
  TARGETS.each do |os, archs|
    archs.each do |arch|
      ENV['GOOS'] = os.to_s
      ENV['GOARCH'] = arch
      system("go build -o build/speedmon_#{os}_#{arch}")
    end
  end
end

task default: :all
