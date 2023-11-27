Name:           tokamakd
Version:        1.0.0
Release:        1%{?dist}
Summary:        Restrict the user resource in HPC login nodes.

License:        GPLv3+
URL:            https://github.com/sonmihpc/tokamak
Source0:       https://github.com/sonmihpc/tokamak/%{name}-%{version}.tar.gz

%description
Restrict the user resource in HPC login nodes.

%prep
%setup -q

%undefine _missing_build_ids_terminate_build
%global debug_package %{nil}

%build

%install
mkdir -p %{buildroot}/%{_sbindir}
mkdir -p %{buildroot}/%{_sysconfdir}/tokamak
mkdir -p %{buildroot}/%{_unitdir}
install -m 0700 tokamakd %{buildroot}/%{_sbindir}/tokamakd
install -m 0644 config.yaml %{buildroot}/%{_sysconfdir}/tokamak/config.yaml
install -m 0644 tokamakd.service %{buildroot}/%{_unitdir}/tokamakd.service

%files
%{_sbindir}/tokamakd
%{_sysconfdir}/tokamak/config.yaml
%{_unitdir}/tokamakd.service

%changelog
* Tue Nov 27 2023 root
- v1.0.0 release