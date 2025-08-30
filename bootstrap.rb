#!/usr/bin/env ruby

require 'net/http'

html_content = Net::HTTP.get(
  URI 'https://www.linuxfromscratch.org/lfs/view/development/chapter03/packages.html'
)

class OpenTag
  attr_accessor :lexeme

  def initialize(lexeme, idx)
    @lexeme = lexeme
    @idx = idx
  end

  @@pattern = /^<[a-zA-Z ]+>/

  def self.match(other)
    @@pattern.match other
  end
end

class CloseTag
  attr_accessor :lexeme

  def initialize(lexeme, idx)
    @lexeme = lexeme
    @idx = idx
  end

  @@pattern = /^<\/[a-zA-Z ]+>/

  def self.match(other)
    @@pattern.match other
  end
end

def scan_html(raw_text)
  len = raw_text.length
  tokens = []
  idx = 0

  while idx < len
    cur = raw_text[idx]
    if cur == ' ' || cur == '\n' || cur == '\t'
      idx += 1
      next
    end

    s = raw_text.slice(idx, len - idx)
    match = OpenTag.match s
    unless match.nil?
      l = match.values_at 0
      tokens.push OpenTag.new(l, idx)
      idx += l.length
      next
    end

    match = CloseTag.match s
    unless match.nil?
      l = match.values_at 0
      tokens.push CloseTag.new(l, idx)
      idx += l.length
      next
    end

    idx += 1 # TODO delete
  end

  tokens
end

def parse_html()
end

tokens = scan_html html_content
p tokens
p tokens.length

#RemoteFile = Struct.new(:remote, :checksum)
#
## https://www.linuxfromscratch.org/lfs/view/development/chapter03/packages.html
#remotes = [
#  RemoteFile.new(
#    'https://ftp.gnu.org/gnu/glibc/glibc-2.41.tar.xz',
#    '19862601af60f73ac69e067d3e9267d4',
#  )
#]
#
#unless Dir.exist? 'ignore'
#  Dir.mkdir 'ignore'
#end
#
#remotes.each do |remote|
#  basename = File.basename remote.remote
#  local_path = "./ignore/#{basename}"
#  unless File.exist? local_path
#    unless system(
#      "curl -L #{remote.remote} -o #{local_path}",
#      exception: true,
#    )
#      raise "Failed to download #{remote.remote}!"
#    end
#  end
#  actual_hash = `md5sum --binary #{local_path} | awk '{print $1}'`.strip
#  if actual_hash != remote.checksum
#    raise "checksum fail! \"#{actual_hash}\" != #{remote.hash}"
#  end
#end
