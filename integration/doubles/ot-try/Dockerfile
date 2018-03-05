# Start from the latest PolarisHQ integrationtester image
FROM 137462835382.dkr.ecr.us-west-1.amazonaws.com/integrationtester

# Copy the source files into place
COPY ./integration/doubles/ot-try/stubby.yml /etc/stubby.yml

CMD ["stubby", "--data", "/etc/stubby.yml", "--stubs", "80", "--tls", "443"]