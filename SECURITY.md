# Security Policy

## Supported Versions

We provide security updates for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 0.3.x   | :white_check_mark: |
| < 0.3   | :x:                |

## Reporting a Vulnerability

We take security vulnerabilities seriously and appreciate your efforts to responsibly disclose your findings. To report a vulnerability, please use one of the following secure methods:

**Please DO NOT open a public GitHub issue for security vulnerabilities.**

### Reporting Methods

1. **GitHub Security Advisories** (Preferred): Use the [Security Advisory](https://github.com/pixelfactoryio/needle/security/advisories/new) feature to privately report vulnerabilities
2. **Email**: Send details to security@pixelfactory.io with "[SECURITY]" in the subject line

### What Constitutes a Vulnerability

A vulnerability is a weakness in the software that could be exploited to:
- Execute arbitrary code or commands
- Bypass authentication or authorization mechanisms
- Access sensitive data or credentials (certificate private keys, configuration data)
- Cause denial of service or system crashes
- Compromise certificate generation or validation
- Bypass ad/tracker blocking mechanisms
- Inject malicious content or certificates
- Compromise the integrity of TLS/DNS handling

### What to Include in Your Report

When reporting a vulnerability, please provide:
- Description of the vulnerability and its potential impact
- Detailed steps to reproduce the issue
- Affected versions
- Any proof-of-concept code or screenshots
- Potential impact assessment
- Suggested fixes or mitigations (if available)

## Coordinated Vulnerability Disclosure

We follow coordinated vulnerability disclosure practices to ensure vulnerabilities are handled responsibly:

### Our Disclosure Timeline

1. **Acknowledgment**: We will acknowledge receipt of your vulnerability report within 48 hours
2. **Initial Assessment**: We will provide an initial assessment and severity rating within 5 business days
3. **Disclosure Coordination**: We request a disclosure embargo of at least 90 days from the initial report to allow time for:
   - Thorough investigation and validation
   - Development and testing of fixes
   - Preparation of security patches and advisories
   - Coordination with downstream users and distributors
4. **Public Disclosure**: After the fix is released, we will publicly disclose the vulnerability within 7 days, or coordinate disclosure timing with the reporter

### Disclosure Process

- We will work with you to understand and validate the vulnerability
- We will keep you informed of our progress throughout the remediation process
- We will credit you in the security advisory (unless you prefer to remain anonymous)
- If a vulnerability is already publicly known or being actively exploited, we may expedite the disclosure timeline to 30 days or less

### Expectations

- **Confidentiality**: Please do not disclose the vulnerability publicly until we have released a fix and agreed upon a disclosure date
- **Good Faith**: We expect security researchers to act in good faith and avoid actions that could harm users, such as accessing unnecessary data, degrading service, or destroying data
- **No Legal Action**: We will not pursue legal action against researchers who discover and report vulnerabilities in accordance with this policy

## Security Update Process

When a vulnerability is confirmed:

1. We will develop and test a security patch
2. We will release the patch in a new version
3. We will publish a security advisory with:
   - Description of the vulnerability
   - Affected versions
   - Fixed versions
   - Workarounds (if available)
   - Credit to the reporter
4. We will notify users through GitHub security advisories and release notes

## Scope

This security policy applies to:
- The needle application and all its components
- Certificate generation and TLS handling
- DNS server integration (CoreDNS)
- HTTP/HTTPS proxy functionality
- Official documentation and examples
- Build and release infrastructure

Out of scope:
- Issues in third-party dependencies (please report to the respective projects)
- Vulnerabilities in applications using needle (unless the root cause is in needle itself)
- Expected behavior (e.g., needle intentionally blocks ad/tracker requests)

## Contact

For security-related questions or concerns:
- **Security Email**: security@pixelfactory.io
- **GitHub**: https://github.com/pixelfactoryio/needle/security/advisories

Thank you for helping keep this project and its users secure!
