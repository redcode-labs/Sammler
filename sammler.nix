{ buildGoModule }:

buildGoModule rec {
  pname = "sammler";
  version = "0.0.1";

  src = builtins.filterSource (path: type: type != "directory" || baseNameOf path != ".git") ./.;

  vendorSha256 = "sha256:hA/wYLQvc3twcacdUhBTD46wjPcIA3DvDGIjphrMIJQ="; 

  subPackages = [ "." ]; 

  runVend = true;

  buildInputs = [ ];
}


