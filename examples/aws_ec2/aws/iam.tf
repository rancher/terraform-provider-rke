# Step 1: Create an IAM role
resource "aws_iam_role" "rke-role" {
  name = "rke-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": {
    "Effect": "Allow",
    "Principal": {"Service": "ec2.amazonaws.com"},
    "Action": "sts:AssumeRole"
  }
}
EOF

}

# Step 2: Add our Access Policy
resource "aws_iam_role_policy" "rke-access-policy" {
name = "rke-access-policy"
role = aws_iam_role.rke-role.id

policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "ec2:Describe*",
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": "ec2:AttachVolume",
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": "ec2:DetachVolume",
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": ["ec2:*"],
      "Resource": ["*"]
    },
    {
      "Effect": "Allow",
      "Action": ["elasticloadbalancing:*"],
      "Resource": ["*"]
    }
  ]
}
EOF

}

# Step 3: Create the Instance Profile
resource "aws_iam_instance_profile" "rke-aws" {
  name = "rke-aws"
  role = aws_iam_role.rke-role.name
}

